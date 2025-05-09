package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
)

var rdb *redis.Client

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortLink() string {
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder

	for i := 0; i < 6; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func InitDB() {
	rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 0})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Eror while connecting to redis: %v", err)
	}
}

func GetFromRedis(key string) (string, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Если Redis не нашел ключ, возвращаем пустую строку
		return "", nil
	}
	if err != nil {
		// Обрабатываем другие ошибки
		return "", err
	}
	return val, nil
}

func SaveToRedis(key string, value string) error {
	ctx := context.Background()
	err := rdb.Set(ctx, key, value, 0).Err() // Сохраняем без времени жизни (TTL)
	if err != nil {
		return err
	}
	return nil
}

// GenerateQRCode генерирует QR код для строки
func GenerateQRCode(shortened string) (string, error) {
	// Ваш код для генерации QR-кода
	// Например, создаем его в виде строки base64
	qr, err := qrcode.New(shortened, qrcode.Medium)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = qr.Write(256, &buf) // генерируем изображение QR-кода в 256px
	if err != nil {
		return "", err
	}

	// Преобразуем изображение в base64 строку
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
