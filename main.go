package main

import (
	"fmt"
	"go-news/pkg/config"
	"log"
	"os"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg := config.Load()

	// Проверяем корректность конфигурации
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Логируем конфигурацию (с замаскированным паролем)
	fmt.Fprintf(os.Stdout, "Application started with config: %s\n", cfg.LogMask())

	// Выводим DSN для дебага (только в разработке)
	if cfg.AppEnv == "development" {
		fmt.Fprintf(os.Stdout, "Database connection string: %s\n", cfg.GetDSN())
	}

	fmt.Fprintf(os.Stdout, "Create data\nListening on port %s\n", cfg.AppPort)
}

