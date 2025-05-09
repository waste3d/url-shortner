package repository

import (
	"url_shorter/db"
	"url_shorter/models"

	"gorm.io/gorm"
)

type LinkRepository struct{}

// Получить ссылку по shortened ID
func (r *LinkRepository) GetLinkByShortened(shortened string) (*models.Link, error) {
	var link models.Link
	if err := db.DB.Where("shortened = ?", shortened).First(&link).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Если не нашли запись, возвращаем nil
		}
		return nil, err
	}
	return &link, nil
}

// Создать ссылку в базе данных
func (r *LinkRepository) CreateLink(link *models.Link) error {
	if err := db.DB.Create(link).Error; err != nil {
		return err
	}
	return nil
}

// Увеличить количество кликов
func (r *LinkRepository) IncrementClicks(shortened string) error {
	if err := db.DB.Model(&models.Link{}).Where("shortened = ?", shortened).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}
