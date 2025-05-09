package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"url_shorter/db"
	"url_shorter/models"
	"url_shorter/repository"
	"url_shorter/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var linkRepo = repository.LinkRepository{}

func CreateLink(ctx *gin.Context) {
	var input models.Link

	// Привязываем JSON к структуре
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка валидности URL
	if !isValidURL(input.Original) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	// Проверка на уникальность shortened ID
	if input.Shortened != "" {
		existingLink, err := linkRepo.GetLinkByShortened(input.Shortened)
		if err == nil && existingLink != nil {
			// Если ссылка уже существует, подтягиваем время создания из БД
			log.Println("Short link already exists, fetching creation time from DB")
			formattedCreatedAt := existingLink.Created_at.Format("02.01.06 15:04")
			formattedExpireAt := existingLink.Expire_at.Format("02.01.06 15:04")

			ctx.JSON(http.StatusOK, gin.H{
				"original":   existingLink.Original,
				"shortened":  existingLink.Shortened,
				"id":         existingLink.ID,
				"clicks":     existingLink.Clicks,
				"created_at": formattedCreatedAt,
				"expire_at":  formattedExpireAt,
			})
			return
		}
	}

	// Если shortened ID не передан, генерируем случайный
	if input.Shortened == "" {
		input.Shortened = utils.GenerateShortLink()
	}

	// Устанавливаем дату создания и дату истечения
	input.Created_at = time.Now()

	if input.Expire_at.IsZero() {
		input.Expire_at = input.Created_at.Add(24 * time.Hour)
	} else if input.Expire_at.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "time must be in the future"})
		return
	}

	// Генерация нового QR-кода
	qrCodeImage, err := utils.GenerateQRCode(input.Original)
	if err != nil {
		log.Println("Ошибка при генерации QR-кода:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Сохранение QR-кода в Redis
	qrCodeKey := fmt.Sprintf("qrcode:%s", input.Shortened)
	err = utils.SaveToRedis(qrCodeKey, qrCodeImage)
	if err != nil {
		log.Println("Ошибка при сохранении QR-кода в Redis:", err)
	}

	// Сохранение ссылки в базе данных через репозиторий
	if err := linkRepo.CreateLink(&input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
		log.Print("db error: ", err.Error())
		return
	}

	// Форматирование времени
	formattedExpireAt := input.Expire_at.Format("02.01.06 15:04")
	formattedCreatedAt := input.Created_at.Format("02.01.06 15:04")

	// Ответ пользователю
	ctx.JSON(http.StatusCreated, gin.H{
		"original":   input.Original,
		"shortened":  input.Shortened,
		"id":         input.ID,
		"clicks":     input.Clicks,
		"created_at": formattedCreatedAt,
		"expire_at":  formattedExpireAt,
		"qr_code":    qrCodeImage,
	})
}

func GetVisitorInfo(ctx *gin.Context) {
	id := ctx.Param("id")

	var link models.Link
	if err := db.DB.Where("id = ?", id).First(&link).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "short link not found"})
		return
	}

	var visitor []models.Visitor
	if err := db.DB.Where("link_id = ?", link.ID).Order("created_at desc").First(&visitor).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No visitors found for this link"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"shortened": link.Shortened, "visitors": visitor})
}

func GetAllLinks(ctx *gin.Context) {
	var links []models.Link

	if err := db.DB.Find(&links).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get links"})
		return
	}
	ctx.JSON(http.StatusOK, links)
}

func GetLinkByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var link models.Link

	if err := db.DB.First(&link, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"original": link.Original, "shortened": link.Shortened, "created_at": link.Created_at, "expire_at": link.Expire_at, "clicks": link.Clicks})
}

func RedirectLink(ctx *gin.Context) {
	shortened := ctx.Param("shortened")

	var link models.Link

	if err := db.DB.Where("shortened = ?", shortened).First(&link).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "short link not found"})
		return
	}

	if time.Now().After(link.Expire_at) {
		ctx.JSON(http.StatusGone, gin.H{"error": "link has expires"})
		return
	}

	visitor := models.Visitor{
		UserIP:    ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
		LinkID:    link.ID,
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&visitor).Error; err != nil {
		log.Println("Error saving visitor: ", err)
	}

	if err := db.DB.Model(&models.Link{}).Where("shortened = ?", shortened).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).Error; err != nil {
		log.Println("Error incrementing clicks:", err)
	}

	ctx.Redirect(http.StatusMovedPermanently, link.Original)
}

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}

// --------------- Работа с QR кодами (вынесена в main.go)----------------
