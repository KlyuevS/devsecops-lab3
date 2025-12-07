# GitHub Secrets Configuration

## Required Secrets for CI/CD Pipeline

Для работы pipeline с конфигурацией базы данных и приложения необходимо установить следующие secrets в GitHub:

### Как добавить secrets:

1. Перейдите в Settings репозитория
2. В левом меню выберите "Secrets and variables" → "Actions"
3. Нажмите "New repository secret"
4. Добавьте каждый из указанных ниже secrets

### Список required secrets:

#### Database Configuration

```
DB_HOST       - Хост базы данных (пример: db.example.com)
DB_PORT       - Порт базы данных (пример: 5432)
DB_USER       - Пользователь БД (пример: postgres)
DB_PASSWORD   - Пароль пользователя БД (ЧУВСТВИТЕЛЬНАЯ ИНФОРМАЦИЯ)
DB_NAME       - Имя базы данных (пример: news_db)
```

#### Application Configuration

```
APP_PORT      - Порт приложения (пример: 8080)
APP_ENV       - Окружение (development/staging/production)
```

### Пример значений для development/staging:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=news_db
APP_PORT=8080
APP_ENV=ci
```

### Пример значений для production:

```
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=<strong-password-here>
DB_NAME=prod_news_db
APP_PORT=8080
APP_ENV=production
```

## Использование в pipeline

Secrets автоматически инжектируются в GitHub Actions workflow через переменные окружения:

```yaml
env:
  DB_HOST: ${{ secrets.DB_HOST || 'localhost' }}
  DB_PORT: ${{ secrets.DB_PORT || '5432' }}
  DB_USER: ${{ secrets.DB_USER || 'postgres' }}
  DB_PASSWORD: ${{ secrets.DB_PASSWORD || 'postgres' }}
  DB_NAME: ${{ secrets.DB_NAME || 'news_db' }}
  APP_PORT: ${{ secrets.APP_PORT || '8080' }}
  APP_ENV: ${{ secrets.APP_ENV || 'ci' }}
```

## Использование в приложении

Конфигурация загружается из переменных окружения в `pkg/config/config.go`:

```go
cfg := config.Load()
// cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName
// cfg.AppPort, cfg.AppEnv
```

## Интеграционные тесты

В `pkg/api/api_integration_test.go` находятся следующие тесты:

- `TestConfigLoadsFromEnvironment` - проверка загрузки переменных из окружения
- `TestConfigUsesDefaultValues` - проверка значений по умолчанию
- `TestConfigDSN` - проверка формирования строки подключения
- `TestConfigValidation` - проверка валидации конфигурации
- `TestConfigLogMask` - проверка маскирования пароля при логировании
- `TestConfigIntegrationWithAPI` - проверка интеграции конфигурации с API

Тесты запускаются автоматически в pipeline и используют переменные окружения из GitHub Secrets.

## Безопасность

⚠️ **ВАЖНО:**

1. **Пароли не должны быть в коде** - используйте только GitHub Secrets
2. **Логи не должны содержать пароли** - используйте `cfg.LogMask()`
3. **Пароли маскируются** при логировании и в outputs
4. **Secrets недоступны** в pull requests с внешних репозиториев
5. **Ротация пароли** - обновляйте secrets периодически

## Локальное тестирование

Для локального тестирования установите переменные окружения:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=news_db
export APP_PORT=8080
export APP_ENV=development

go test -v ./pkg/api
go run main.go
```
