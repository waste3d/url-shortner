package models

import (
	"time"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	ID         int       `json:"id"`
	Original   string    `json:"original"`
	Shortened  string    `json:"shortened"`
	Created_at time.Time `json:"created_at"`
	Expire_at  time.Time `json:"expire_at"`
	Clicks     int       `json:"clicks"`
	ShortID    string    `json:"short_id"`
}

type Visitor struct {
	gorm.Model
	ID        int       `json:"id"`
	UserAgent string    `json:"user_agent"`
	UserIP    string    `json:"user_IP"`
	LinkID    int       `json:"link_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Users struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
