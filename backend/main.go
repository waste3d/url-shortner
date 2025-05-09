package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"time"
	"url_shorter/db"
	controllers "url_shorter/handlers"
	"url_shorter/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
)

var (
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx = context.Background()
)

func main() {

	utils.InitDB()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/links", controllers.GetAllLinks)
	r.POST("/links", controllers.CreateLink)
	r.GET("/:shortened", controllers.RedirectLink)
	r.GET("/links/:id", controllers.GetLinkByID)
	r.GET("/links/:id/visitors", controllers.GetVisitorInfo)
	r.GET("/qr/view", generateQRHandler)
	r.GET("qr/download", downloadQRHandler)

	r.Run(":8080")
}

func generateQRHandler(c *gin.Context) {
	url := c.Query("url")
	hexColor := c.Query("color")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	qrBytes, err := GetCachedQR(url, hexColor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate QR"})
		return
	}

	c.Data(http.StatusOK, "image/png", qrBytes)
}

// скачивание QR
func downloadQRHandler(c *gin.Context) {
	url := c.Query("url")
	hexColor := c.Query("color")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	qrBytes, err := GetCachedQR(url, hexColor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate QR"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=qr.png")
	c.Data(http.StatusOK, "image/png", qrBytes)
}

// генерация QR с возможным цветом
func GenerateQRBytes(url string, hexColor string) ([]byte, error) {
	col := &color.RGBA{0, 0, 0, 255} // default black
	if hexColor != "" {
		if parsedColor, err := parseHexColor(hexColor); err == nil {
			col = parsedColor
		}
	}

	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qr.ForegroundColor = col

	var buf bytes.Buffer
	if err := png.Encode(&buf, qr.Image(256)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// из #000000 в RGB
func parseHexColor(s string) (*color.RGBA, error) {
	if len(s) != 6 {
		return nil, fmt.Errorf("invalid length")
	}
	var r, g, b uint8
	_, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return nil, err
	}
	return &color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

func GetCachedQR(url, hexColor string) ([]byte, error) {
	key := GenerateCacheKey(url, hexColor)
	ctx := context.Background()

	// попытка достать из redis

	// QR не найден - генерируем
	qrBytes, err := GenerateQRBytes(url, hexColor)
	if err != nil {
		return nil, err
	}

	// сохраняем в кеш
	err = rdb.Set(ctx, key, qrBytes, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Ошибка при сохранении в Redis: %v", err)
	}

	return qrBytes, nil
}

func GenerateCacheKey(url, color string) string {
	hash := md5.Sum([]byte(url + color))
	return "qr:" + hex.EncodeToString(hash[:])
}
