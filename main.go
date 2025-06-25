package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"workmateTestProject/internal/model"
	"workmateTestProject/internal/service"
	"workmateTestProject/internal/storage"
)

// errorResponse отвечает JSON с сообщением об ошибке
func errorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	type errT struct {
		Error string `json:"error"`
	}
	// Обрабатываем возможную ошибку кодирования JSON, чтобы не оставлять её без внимания
	if err := json.NewEncoder(w).Encode(errT{Error: message}); err != nil {
		log.Printf("Ошибка при кодировании JSON-ошибки: %v", err)
	}
}

// loggingMiddleware логирует входящие запросы и время их обработки
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Начало %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Завершено %s %s за %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	// Создаём in-memory хранилище задач
	store := storage.NewInMemoryTaskStore()

	// Настраиваем маршрутизатор
	r := mux.NewRouter()

	// подмешиваем логирование
	h := loggingMiddleware(r)

	// Роуты для работы с задачами
	r.HandleFunc("/tasks", createTaskHandler(store)).Methods(http.MethodPost)
	r.HandleFunc("/tasks", listTasksHandler(store)).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}", getTaskHandler(store)).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}", deleteTaskHandler(store)).Methods(http.MethodDelete)

	// Определяем порт из переменной окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Создаём HTTP сервер
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Сервер запущен на :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Получен сигнал завершения, выключаем сервер...")

	// Пытаемся корректно завершить с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Сервер не завершился корректно: %v", err)
	}
	log.Println("Сервер завершён")
}

// createTaskHandler обрабатывает создание новой задачи
func createTaskHandler(store storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Создаём новую задачу
		id := uuid.New()
		task := &model.Task{
			ID:        id,
			Status:    model.StatusPending,
			CreatedAt: time.Now(),
		}
		store.Create(task)
		// Запускаем обработку задачи
		service.StartProcessing(task)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(task); err != nil {
			errorResponse(w, http.StatusInternalServerError, "Ошибка кодирования ответа")
		}
	}
}

// getTaskHandler возвращает информацию о задаче по ID
func getTaskHandler(store storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := uuid.Parse(vars["id"])
		if err != nil {
			errorResponse(w, http.StatusBadRequest, "Неверный UUID")
			return
		}
		task, ok := store.Get(id)
		if !ok {
			errorResponse(w, http.StatusNotFound, "Задача не найдена")
			return
		}

		// Подготавливаем ответ с вычислением длительности
		type respT struct {
			ID         uuid.UUID        `json:"id"`
			Status     model.TaskStatus `json:"status"`
			CreatedAt  time.Time        `json:"created_at"`
			StartedAt  *time.Time       `json:"started_at,omitempty"`
			FinishedAt *time.Time       `json:"finished_at,omitempty"`
			Duration   *string          `json:"duration,omitempty"`
			Result     string           `json:"result,omitempty"`
			Error      string           `json:"error,omitempty"`
		}
		resp := respT{
			ID:         task.ID,
			Status:     task.Status,
			CreatedAt:  task.CreatedAt,
			StartedAt:  task.StartedAt,
			FinishedAt: task.FinishedAt,
			Result:     task.Result,
			Error:      task.Error,
		}
		if task.StartedAt != nil && task.FinishedAt != nil {
			d := task.FinishedAt.Sub(*task.StartedAt).String()
			resp.Duration = &d
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			errorResponse(w, http.StatusInternalServerError, "Ошибка кодирования ответа")
		}
	}
}

// listTasksHandler возвращает список всех задач
func listTasksHandler(store storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks := store.List()

		// Формируем ответ по аналогии с getTaskHandler
		type respT struct {
			ID         uuid.UUID        `json:"id"`
			Status     model.TaskStatus `json:"status"`
			CreatedAt  time.Time        `json:"created_at"`
			StartedAt  *time.Time       `json:"started_at,omitempty"`
			FinishedAt *time.Time       `json:"finished_at,omitempty"`
			Duration   *string          `json:"duration,omitempty"`
			Result     string           `json:"result,omitempty"`
			Error      string           `json:"error,omitempty"`
		}
		responses := make([]respT, 0, len(tasks))
		for _, task := range tasks {
			r := respT{
				ID:         task.ID,
				Status:     task.Status,
				CreatedAt:  task.CreatedAt,
				StartedAt:  task.StartedAt,
				FinishedAt: task.FinishedAt,
				Result:     task.Result,
				Error:      task.Error,
			}
			if task.StartedAt != nil && task.FinishedAt != nil {
				d := task.FinishedAt.Sub(*task.StartedAt).String()
				r.Duration = &d
			}
			responses = append(responses, r)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			errorResponse(w, http.StatusInternalServerError, "Ошибка кодирования ответа")
		}
	}
}

// deleteTaskHandler удаляет задачу по ID
func deleteTaskHandler(store storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := uuid.Parse(vars["id"])
		if err != nil {
			errorResponse(w, http.StatusBadRequest, "Неверный UUID")
			return
		}
		store.Delete(id)
		w.WriteHeader(http.StatusNoContent)
	}
}
