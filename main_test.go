package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
	"workmateTestProject/internal/model"
	"workmateTestProject/internal/service"
	"workmateTestProject/internal/storage"
)

// setupRouter создаёт новый роутер с зарегистрированными хендлерами,
// используется для тестирования HTTP API без запуска реального сервера.
func setupRouter() http.Handler {
	store := storage.NewInMemoryTaskStore()
	// Подменяем симулятор на мгновенный для ускорения тестов
	service.SimulateWorkFunc = func() (string, error) {
		return "fast-result", nil
	}

	r := mux.NewRouter()
	r.HandleFunc("/tasks", createTaskHandler(store)).Methods(http.MethodPost)
	r.HandleFunc("/tasks", listTasksHandler(store)).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}", getTaskHandler(store)).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}", deleteTaskHandler(store)).Methods(http.MethodDelete)

	// Логирование не требуется в тестах, возвращаем роутер напрямую
	return r
}

// TestCreateAndGetAndDelete проверяет сценарий создания, получения и удаления задачи через HTTP API.
func TestCreateAndGetAndDelete(t *testing.T) {
	h := setupRouter()

	// 1. Создаём задачу через POST /tasks
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("ожидался код 201 Created, получили %d", rec.Code)
	}

	// Декодируем ответ JSON с полем id
	var created model.Task
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}

	// Проверяем, что ID валидный UUID и статус Pending
	if _, err := uuid.Parse(created.ID.String()); err != nil {
		t.Errorf("поле id не является корректным UUID: %v", err)
	}
	if created.Status != model.StatusPending {
		t.Errorf("ожидался статус Pending, получили %s", created.Status)
	}

	// 2. Сразу запрашиваем GET /tasks/{id}, ожидаем Pending (обработка мгновенная)
	rec = httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/tasks/"+created.ID.String(), nil)
	h.ServeHTTP(rec, getReq)
	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался код 200 OK, получили %d", rec.Code)
	}
	var fetched struct {
		ID     uuid.UUID        `json:"id"`
		Status model.TaskStatus `json:"status"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&fetched); err != nil {
		t.Fatalf("ошибка декодирования GET response: %v", err)
	}
	if fetched.Status != model.StatusCompleted {
		t.Errorf("ожидался статус Completed благодаря fast-result, получили %s", fetched.Status)
	}

	// 3. Удаляем задачу через DELETE /tasks/{id}
	rec = httptest.NewRecorder()
	delReq := httptest.NewRequest(http.MethodDelete, "/tasks/"+created.ID.String(), nil)
	h.ServeHTTP(rec, delReq)
	if rec.Code != http.StatusNoContent {
		t.Errorf("ожидался код 204 No Content при удалении, получили %d", rec.Code)
	}

	// 4. После удаления GET снова должен вернуть 404
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, getReq)
	if rec.Code != http.StatusNotFound {
		t.Errorf("ожидался код 404 Not Found после удаления, получили %d", rec.Code)
	}
}
