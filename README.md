


Система использует следующую модель данных для задач:

```go
type Task struct {
    ID              int    // ID задачи
    ResponsibleID   int    // id ответственного
    ResponsibleName string // имя ответственного
    DueDate         int64  // срок выполнения задачи (Unix timestamp)
    AssignedAt      int64  // срок постановки задачи / дата назначения (Unix timestamp)
    Context         string // контекст / описание задачи
}
```



## API Endpoints
- GET /posts - получение всех задач
- POST /posts - создание новой задачи
- PUT /posts - обновление существующей задачи
- DELETE /posts - удаление задачи

## Тестирование
В проекте реализован тестовый клиент (`cmd/test/test_api.go`), который проверяет все  операции:
- GET - получение списка всех задач
- POST - создание новой тестовой задачи
- PUT - обновление созданной задачи
- DELETE - удаление тестовой задачи



## Запуск проверок
1. Создать тестовую запись:
```bash
curl -k -X POST https://localhost/posts \
  -H "Content-Type: application/json" \
  -d '{
        "id": 2,
        "responsible_id": 101,
        "responsible_name": "SergeyKlyuev",
        "context": "DevOps cool!",
        "assigned_at": 4444,
        "due_date": 888
      }'

```

2. Получить все записи:
```bash
curl -k https://localhost/posts
```

3. Посмотреть записи в БД:
```bash
docker exec -it news_app-db-1 \
  psql -U news_user -d news -c "SELECT * FROM posts;"
```
