package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-news/pkg/storage"
	"io"
	"net/http"
)

const (
	greenColor = "\033[32m"
	redColor   = "\033[31m"
	resetColor = "\033[0m"
)

func printResult(operation string, err error) {
	if err != nil {
		fmt.Printf("%s%s: ERROR - %v%s\n", redColor, operation, err, resetColor)
		return
	}
	fmt.Printf("%s%s: SUCCESS%s\n", greenColor, operation, resetColor)
}

func main() {
	baseURL := "http://localhost:8080/posts"

	// GET - получение всех задач
	resp, err := http.Get(baseURL)
	if err != nil {
		printResult("GET", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Current tasks: %s\n", string(body))
	printResult("GET", nil)

	// POST - создание новой задачи
	task := storage.Task{
		ID:              100,
		ResponsibleID:   1,
		ResponsibleName: "Test User",
		Context:         "Test Task",
		AssignedAt:      1673891100,
		DueDate:         1674064800,
	}

	jsonData, _ := json.Marshal(task)
	resp, err = http.Post(baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		printResult("POST", err)
		return
	}
	resp.Body.Close()
	printResult("POST", nil)

	// PUT - обновление задачи
	task.Context = "Updated Test Task"
	jsonData, _ = json.Marshal(task)
	req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		printResult("PUT", err)
		return
	}
	resp.Body.Close()
	printResult("PUT", nil)

	// DELETE - удаление задачи
	jsonData, _ = json.Marshal(task)
	req, _ = http.NewRequest(http.MethodDelete, baseURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		printResult("DELETE", err)
		return
	}
	resp.Body.Close()
	printResult("DELETE", nil)
}
