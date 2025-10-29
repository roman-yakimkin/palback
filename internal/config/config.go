package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string

	FrontendOrigin string

	MinIOEndpoint          string
	MinIOAccessKey         string
	MinIOSecretKey         string
	MinIOBucketMain        string
	MinIOBucketUserAvatars string

	RedisAddr      string
	RedisUsername  string
	RedisPassword  string
	RedisSecretKey string

	SessionDays int

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден")
	}

	if err := godotenv.Load(".env.local"); err == nil {
		log.Println("Загружен .env.local")
	} else if !os.IsNotExist(err) {
		log.Printf("Ошибка при загрузке .env.local: %v", err)
	}

	cfg := &Config{
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "mydb"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),

		MinIOEndpoint:          getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinIOAccessKey:         getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey:         getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucketMain:        getEnv("MINIO_BUCKET_MAIN", "main"),
		MinIOBucketUserAvatars: getEnv("MINIO_BUCKET_USER_AVATARS", "user-avatars"),

		RedisAddr:      getEnv("REDIS_ADDR", "redis:6379"),
		RedisUsername:  getEnv("REDIS_USERNAME", ""),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisSecretKey: getEnv("REDIS_SECRET_KEY", "secret-key-32-bytes-long-12345678"),

		SMTPHost:     getEnv("SMTP_HOST", "localhost"),
		SMTPPort:     getEnv("SMTP_PORT", "1025"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "no-reply@palomniki.su"),
	}

	cfg.SessionDays, err = strconv.Atoi(getEnv("SESSION_DAYS", "7"))
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetLang() string {
	return "ru"
}
