package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"palback/internal/infra/storage"
	"strconv"

	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"palback/internal/config"
	handler "palback/internal/delivery/http"
	"palback/internal/infra/repository"
	"palback/internal/usecase"
)

func main() {

	// Загрузка данных из конфига
	cfg := config.Load()

	// Подключение БД
	dbConnString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	db, err := sql.Open(cfg.DBDriver, dbConnString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Подключение minio
	minioClient, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: false, // отключено, т.к. MinIO в Docker без TLS
	})
	if err != nil {
		log.Fatal("Ошибка инициализации MinIO:", err)
	}

	// Подключение redis
	rediStore, err := redistore.NewRediStore(
		10,
		"tcp",
		cfg.RedisAddr,
		cfg.RedisUsername,
		cfg.RedisPassword,
		[]byte(cfg.RedisSecretKey),
	)
	if err != nil {
		log.Fatal("Ошибка инициализации Redis:", err)
	}
	defer rediStore.Close()

	days, err := strconv.Atoi(cfg.SessionDays)
	if err != nil {
		log.Fatal("Ошибка получения длительности сессии", err)
	}
	rediStore.SetMaxAge(86400 * days)

	var _ sessions.Store = rediStore

	// Инициализация слоёв приложения
	_ = storage.NewMinioStorage(minioClient, cfg.MinIOBucketUserAvatars)

	countryRepo := repository.NewCountryRepo(db)
	countryService := usecase.NewCountryUseCase(countryRepo)
	countryHandler := handler.NewCountryHandler(countryService)

	regionRepo := repository.NewRegionRepo(db)
	regionService := usecase.NewRegionUseCase(countryService, regionRepo)
	regionHandler := handler.NewRegionHandler(regionService)

	cityTypeRepo := repository.NewCityTypeRepo(db)
	cityTypeService := usecase.NewCityTypeUseCase(cityTypeRepo)
	cityTypeHandler := handler.NewCityTypeHandler(cityTypeService)

	placeTypeRepo := repository.NewPlaceTypeRepo(db)
	placeTypeService := usecase.NewPlaceTypeUseCase(placeTypeRepo)
	placeTypeHandler := handler.NewPlaceTypeHandler(placeTypeService)

	// Инициализация рутера
	router := handler.NewRouter(
		cfg,
		countryHandler,
		regionHandler,
		cityTypeHandler,
		placeTypeHandler,
	)

	if err := router.Start(":" + cfg.ServerPort); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
