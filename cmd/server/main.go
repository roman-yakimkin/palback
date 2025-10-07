package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"palback/internal/config"
	handler "palback/internal/delivery/http"
	"palback/internal/repository"
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

	// Инициализация слоёв приложения
	countryRepo := repository.NewCountryRepo(db)
	countryService := usecase.NewCountryUseCase(countryRepo)
	countryHandler := handler.NewCountryHandler(countryService)

	regionRepo := repository.NewRegionRepo(db)
	regionService := usecase.NewRegionUseCase(countryService, regionRepo)
	regionHandler := handler.NewRegionHandler(regionService)

	cityTypeRepo := repository.NewCityTypeRepo(db)
	cityTypeService := usecase.NewCityTypeUseCase(cityTypeRepo)
	cityTypeHandler := handler.NewCityTypeHandler(cityTypeService)

	// Инициализация рутера
	router := handler.NewRouter(
		countryHandler,
		regionHandler,
		cityTypeHandler,
	)

	if err := router.Start(":" + cfg.ServerPort); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
