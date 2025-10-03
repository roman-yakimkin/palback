package main

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	handler "palback/internal/delivery/http"
	"palback/internal/repository"
	"palback/internal/usecase"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5444 user=paldev password=paldev dbname=paldev sslmode=disable")
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

	router := handler.NewRouter(countryHandler)

	if err := router.Start(":8080"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
