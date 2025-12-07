package api

import (
	"bytes"
	"encoding/json"
	"go-news/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockDB - mock реализация хранилища для тестирования
type MockDB struct {
	tasks []storage.Task
}

func (m *MockDB) Tasks() ([]storage.Task, error) {
	return m.tasks, nil
}

func (m *MockDB) AddTask(task storage.Task) error {
	m.tasks = append(m.tasks, task)
	return nil
}

func (m *MockDB) UpdateTask(task storage.Task) error {
	for i, t := range m.tasks {
		if t.ID == task.ID {
			m.tasks[i] = task
			return nil
		}
	}
	return nil
}

func (m *MockDB) DeleteTask(task storage.Task) error {
	for i, t := range m.tasks {
		if t.ID == task.ID {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return nil
}

// Test 1: GET /posts - получение всех задач
func TestGetPosts(t *testing.T) {
	mockDB := &MockDB{
		tasks: []storage.Task{
			{
				ID:              1,
				ResponsibleID:   1,
				ResponsibleName: "John Doe",
				Context:         "Task 1",
				AssignedAt:      1673891100,
				DueDate:         1674064800,
			},
			{
				ID:              2,
				ResponsibleID:   2,
				ResponsibleName: "Jane Smith",
				Context:         "Task 2",
				AssignedAt:      1673891200,
				DueDate:         1674064900,
			},
		},
	}

	api := New(mockDB)

	// Создаем HTTP запрос
	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	// Вызываем handler через router
	api.Router().ServeHTTP(w, req)

	// Проверяем статус код
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем тело ответа
	var posts []storage.Task
	err := json.NewDecoder(w.Body).Decode(&posts)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Проверяем что получили 2 задачи
	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}

	// Проверяем первую задачу
	if posts[0].ID != 1 || posts[0].ResponsibleName != "John Doe" {
		t.Errorf("First post data mismatch: ID=%d, Name=%s", posts[0].ID, posts[0].ResponsibleName)
	}
}

// Test 2: POST /posts - добавление новой задачи
func TestAddPost(t *testing.T) {
	mockDB := &MockDB{tasks: []storage.Task{}}
	api := New(mockDB)

	newTask := storage.Task{
		ID:              1,
		ResponsibleID:   1,
		ResponsibleName: "Alice Johnson",
		Context:         "New Task",
		AssignedAt:      1673891100,
		DueDate:         1674064800,
	}

	jsonData, _ := json.Marshal(newTask)
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, req)

	// Проверяем статус код
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем что задача была добавлена в БД
	if len(mockDB.tasks) != 1 {
		t.Errorf("Expected 1 task in database, got %d", len(mockDB.tasks))
	}

	// Проверяем параметры добавленной задачи
	if mockDB.tasks[0].ResponsibleName != "Alice Johnson" {
		t.Errorf("Expected responsible name 'Alice Johnson', got '%s'", mockDB.tasks[0].ResponsibleName)
	}
}

// Test 3: PUT /posts - обновление задачи
func TestUpdatePost(t *testing.T) {
	mockDB := &MockDB{
		tasks: []storage.Task{
			{
				ID:              1,
				ResponsibleID:   1,
				ResponsibleName: "Old Name",
				Context:         "Old Task",
				AssignedAt:      1673891100,
				DueDate:         1674064800,
			},
		},
	}

	api := New(mockDB)

	updatedTask := storage.Task{
		ID:              1,
		ResponsibleID:   1,
		ResponsibleName: "New Name",
		Context:         "Updated Task",
		AssignedAt:      1673891100,
		DueDate:         1674064800,
	}

	jsonData, _ := json.Marshal(updatedTask)
	req := httptest.NewRequest(http.MethodPut, "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, req)

	// Проверяем статус код
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем что задача была обновлена в БД
	if len(mockDB.tasks) != 1 {
		t.Errorf("Expected 1 task in database, got %d", len(mockDB.tasks))
	}

	// Проверяем что данные обновлены
	if mockDB.tasks[0].ResponsibleName != "New Name" {
		t.Errorf("Expected responsible name 'New Name', got '%s'", mockDB.tasks[0].ResponsibleName)
	}

	if mockDB.tasks[0].Context != "Updated Task" {
		t.Errorf("Expected context 'Updated Task', got '%s'", mockDB.tasks[0].Context)
	}
}

// Test 4: DELETE /posts - удаление задачи
func TestDeletePost(t *testing.T) {
	mockDB := &MockDB{
		tasks: []storage.Task{
			{
				ID:              1,
				ResponsibleID:   1,
				ResponsibleName: "John Doe",
				Context:         "Task 1",
				AssignedAt:      1673891100,
				DueDate:         1674064800,
			},
			{
				ID:              2,
				ResponsibleID:   2,
				ResponsibleName: "Jane Smith",
				Context:         "Task 2",
				AssignedAt:      1673891200,
				DueDate:         1674064900,
			},
		},
	}

	api := New(mockDB)

	taskToDelete := storage.Task{
		ID: 1,
	}

	jsonData, _ := json.Marshal(taskToDelete)
	req := httptest.NewRequest(http.MethodDelete, "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, req)

	// Проверяем статус код
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Проверяем что задача была удалена из БД
	if len(mockDB.tasks) != 1 {
		t.Errorf("Expected 1 task remaining in database, got %d", len(mockDB.tasks))
	}

	// Проверяем что остается правильная задача
	if mockDB.tasks[0].ID != 2 {
		t.Errorf("Expected remaining task with ID 2, got %d", mockDB.tasks[0].ID)
	}
}
