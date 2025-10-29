package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"palback/internal/config"
	handler "palback/internal/delivery/http"
	"palback/internal/infra/email"
	"palback/internal/infra/rate"
	"palback/internal/infra/repository"
	"palback/internal/infra/session"
	"palback/internal/infra/storage"
	"palback/internal/usecase"
)

func main() {
	// Загрузка данных из конфига
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации", err)
	}

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

	// Инициализация redis
	redisPool := initRedisPool(cfg)
	defer redisPool.Close()

	// Инициализация слоёв приложения
	_ = storage.NewMinioStorage(minioClient, cfg.MinIOBucketUserAvatars)

	redisStorage := storage.NewRedisStorage(redisPool)

	mailSender := email.NewSMTPSender(cfg)

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

	roleRepo := repository.NewRoleRepo()
	roleService := usecase.NewRoleUseCase(roleRepo)

	userRepo := repository.NewUserRepo(db)

	auth, err := session.NewRedigoAuthenticator(cfg, redisPool, userRepo)
	if err != nil {
		log.Fatal("Ошибка инициализации аутентификатора", err)
	}

	rateLimiter := rate.NewRedigoRateLimiter(redisPool)

	userService := usecase.NewUserUseCase(roleService, mailSender, redisStorage, userRepo)
	userHandler := handler.NewUserHandler(userService, auth, rateLimiter)

	// Инициализация рутера
	router := handler.NewRouter(
		cfg,
		auth,
		rateLimiter,
		countryHandler,
		regionHandler,
		cityTypeHandler,
		placeTypeHandler,
		userHandler,
	)

	if err := router.Start(":" + cfg.ServerPort); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func initRedisPool(cfg *config.Config) *redis.Pool {
	return &redis.Pool{
		// Максимальное количество "простаивающих" (idle) соединений в пуле
		MaxIdle: 10,

		// Максимальное количество активных соединений (включая idle)
		// 0 = без ограничений (не рекомендуется)
		MaxActive: 50,

		// Время, в течение которого простаивающее соединение может оставаться в пуле
		IdleTimeout: 240 * time.Second, // 4 минуты

		// Функция подключения к Redis
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				cfg.RedisAddr, // адрес вашего Redis
				// Опционально: таймауты
				redis.DialConnectTimeout(5*time.Second),
				redis.DialReadTimeout(3*time.Second),
				redis.DialWriteTimeout(3*time.Second),
			)
			if err != nil {
				return nil, err
			}

			// Опционально: аутентификация (если Redis защищён паролем)
			// if _, err := conn.Do("AUTH", "your-redis-password"); err != nil {
			//     conn.Close()
			//     return nil, err
			// }

			// Опционально: выбор базы (по умолчанию — 0)
			// if _, err := conn.Do("SELECT", 1); err != nil {
			//     conn.Close()
			//     return nil, err
			// }

			return conn, nil
		},

		// Проверка соединения перед использованием из пула
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			// Проверяем соединение каждые 10 секунд
			if time.Since(t) < 10*time.Second {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}
