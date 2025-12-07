package main

import (
	"log"
	"net/http"

	"go-news/pkg/api"
	"go-news/pkg/storage"
	"go-news/pkg/storage/postgres"
)

type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Подключение к Postgres
	connStr := "postgres://news_user:news_pass@db:5432/news"

	db, err := postgres.New(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	srv := server{
		db: db,
	}

	// Создаём API с подключением к БД
	srv.api = api.New(srv.db)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", srv.api.Router()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
