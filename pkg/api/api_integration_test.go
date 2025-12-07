package api

import (
	"go-news/pkg/config"
	"os"
	"testing"
)

// TestConfigLoadsFromEnvironment проверяет загрузку конфигурации из переменных окружения
func TestConfigLoadsFromEnvironment(t *testing.T) {
	// Установим переменные окружения
	os.Setenv("DB_HOST", "test-host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "test-user")
	os.Setenv("DB_PASSWORD", "test-pass")
	os.Setenv("DB_NAME", "test-db")
	os.Setenv("APP_PORT", "9000")
	os.Setenv("APP_ENV", "test")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("APP_ENV")
	}()

	cfg := config.Load()

	// Проверяем загруженные значения
	if cfg.DBHost != "test-host" {
		t.Errorf("Expected DB_HOST='test-host', got '%s'", cfg.DBHost)
	}
	if cfg.DBPort != "5433" {
		t.Errorf("Expected DB_PORT='5433', got '%s'", cfg.DBPort)
	}
	if cfg.DBUser != "test-user" {
		t.Errorf("Expected DB_USER='test-user', got '%s'", cfg.DBUser)
	}
	if cfg.DBPassword != "test-pass" {
		t.Errorf("Expected DB_PASSWORD='test-pass', got '%s'", cfg.DBPassword)
	}
	if cfg.DBName != "test-db" {
		t.Errorf("Expected DB_NAME='test-db', got '%s'", cfg.DBName)
	}
	if cfg.AppPort != "9000" {
		t.Errorf("Expected APP_PORT='9000', got '%s'", cfg.AppPort)
	}
	if cfg.AppEnv != "test" {
		t.Errorf("Expected APP_ENV='test', got '%s'", cfg.AppEnv)
	}
}

// TestConfigUsesDefaultValues проверяет использование значений по умолчанию
func TestConfigUsesDefaultValues(t *testing.T) {
	// Очищаем переменные окружения
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("APP_ENV")

	cfg := config.Load()

	// Проверяем значения по умолчанию
	if cfg.DBHost != "localhost" {
		t.Errorf("Expected DB_HOST='localhost', got '%s'", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("Expected DB_PORT='5432', got '%s'", cfg.DBPort)
	}
	if cfg.AppPort != "8080" {
		t.Errorf("Expected APP_PORT='8080', got '%s'", cfg.AppPort)
	}
}

// TestConfigDSN проверяет формирование строки подключения к БД
func TestConfigDSN(t *testing.T) {
	os.Setenv("DB_HOST", "myhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "myuser")
	os.Setenv("DB_PASSWORD", "mypass")
	os.Setenv("DB_NAME", "mydb")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()

	cfg := config.Load()
	dsn := cfg.GetDSN()

	expected := "postgres://myuser:mypass@myhost:5433/mydb?sslmode=disable"
	if dsn != expected {
		t.Errorf("Expected DSN='%s', got '%s'", expected, dsn)
	}
}

// TestConfigValidation проверяет валидацию конфигурации
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "Valid config",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
				"DB_USER": "user",
				"DB_NAME": "db",
				"APP_PORT": "8080",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Установим переменные
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg := config.Load()
			err := cfg.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Очистка
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

// TestConfigLogMask проверяет маскирование пароля при логировании
func TestConfigLogMask(t *testing.T) {
	os.Setenv("DB_PASSWORD", "secret-password")
	defer os.Unsetenv("DB_PASSWORD")

	cfg := config.Load()
	logOutput := cfg.LogMask()

	// Проверяем, что пароль замаскирован
	if logOutput == "" {
		t.Error("LogMask() returned empty string")
	}

	// Проверяем, что реальный пароль не содержится в логе
	if contains(logOutput, "secret-password") {
		t.Error("LogMask() exposed the actual password")
	}

	// Проверяем, что маска присутствует
	if !contains(logOutput, "***") {
		t.Error("LogMask() does not contain password mask")
	}
}

// TestConfigIntegrationWithAPI проверяет работу конфигурации с API
func TestConfigIntegrationWithAPI(t *testing.T) {
	os.Setenv("DB_HOST", "integration-test-host")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "test-user")
	os.Setenv("DB_PASSWORD", "test-pass")
	os.Setenv("DB_NAME", "test-db")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()

	// Загружаем конфигурацию
	cfg := config.Load()

	// Проверяем, что конфигурация валидна
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}

	// Проверяем, что DSN может быть сформирован
	dsn := cfg.GetDSN()
	if dsn == "" {
		t.Error("Failed to generate DSN from config")
	}

	// Проверяем наличие host в DSN
	if !contains(dsn, "integration-test-host") {
		t.Error("DSN does not contain configured host")
	}

	// Проверяем наличие database name в DSN
	if !contains(dsn, "test-db") {
		t.Error("DSN does not contain configured database name")
	}
}

// Вспомогательная функция для проверки наличия строки
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
