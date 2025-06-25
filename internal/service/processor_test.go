package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"workmateTestProject/internal/model"
)

// TestStartProcessing_Success проверяет корректное обновление полей Task при успешном завершении обработки.
// Мы заменяем simulateWorkFunc на быстрый заглушечный вариант, чтобы тест шел мгновенно.
func TestStartProcessing_Success(t *testing.T) {
	orig := SimulateWorkFunc
	defer func() { SimulateWorkFunc = orig }()

	// Подменяем функцию работы на мгновенную
	SimulateWorkFunc = func() (string, error) {
		// Возвращаем заранее известный результат
		return "stub-result", nil
	}

	// Создаем новую задачу с неопределенным статусом
	id := uuid.New()
	task := &model.Task{ID: id, Status: model.StatusPending}

	// Запускаем обработку
	StartProcessing(task)

	// Ожидаем выполнения горутины (макс 1 сек)
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if task.Status != model.StatusPending {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Проверяем, что статус стал Completed
	if task.Status != model.StatusCompleted {
		t.Fatalf("ожидался статус Completed, получили %v", task.Status)
	}

	// Проверяем, что Result установлен из stub
	if task.Result != "stub-result" {
		t.Errorf("ожидался Result stub-result, получили %v", task.Result)
	}

	// Проверяем, что поля StartedAt и FinishedAt заполнены и FinishedAt позже StartedAt
	if task.StartedAt == nil || task.FinishedAt == nil {
		t.Fatal("Ожидалось, что StartedAt и FinishedAt не nil после обработки")
	}
	if task.FinishedAt.Before(*task.StartedAt) {
		t.Errorf("FinishedAt должно быть позже StartedAt, got %v before %v", task.FinishedAt, task.StartedAt)
	}
}

// TestStartProcessing_Failure проверяет поведение при ошибке simulateWorkFunc.
func TestStartProcessing_Failure(t *testing.T) {
	orig := SimulateWorkFunc
	defer func() { SimulateWorkFunc = orig }()

	// Подменяем функцию на возвращающую ошибку
	SimulateWorkFunc = func() (string, error) {
		return "", fmt.Errorf("simulated error")
	}

	// Создаем задачу и запускаем обработку
	id := uuid.New()
	task := &model.Task{ID: id, Status: model.StatusPending}
	StartProcessing(task)

	// Ждем завершения
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if task.Status != model.StatusPending {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Ожидаем статус Failed
	if task.Status != model.StatusFailed {
		t.Fatalf("ожидался статус Failed, получили %v", task.Status)
	}

	// Проверяем, что Error содержит наш текст
	if task.Error == "" || task.Error != "simulated error" {
		t.Errorf("ожидалась ошибка simulated error, получили %v", task.Error)
	}
}
