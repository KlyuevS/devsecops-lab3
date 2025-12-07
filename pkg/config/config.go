package config

import (
	"fmt"
	"os"
)

// Config хранит конфигурацию приложения
type Config struct {
	// Настройки базы данных
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Настройки приложения
	AppPort string
	AppEnv  string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	cfg := &Config{
		// База данных
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "news_db"),

		// Приложение
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),
	}

	return cfg
}

// getEnv получает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDSN возвращает строку подключения к PostgreSQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBPort == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.AppPort == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	return nil
}

// LogMask возвращает конфигурацию с замаскированным паролем для логирования
func (c *Config) LogMask() string {
	maskedPassword := "***"
	if c.DBPassword != "" {
		maskedPassword = "***"
	}
	return fmt.Sprintf(
		"Config{DBHost: %s, DBPort: %s, DBUser: %s, DBPassword: %s, DBName: %s, AppPort: %s, AppEnv: %s}",
		c.DBHost, c.DBPort, c.DBUser, maskedPassword, c.DBName, c.AppPort, c.AppEnv,
	)
}
