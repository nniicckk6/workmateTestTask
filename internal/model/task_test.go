package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestTaskJSON_MarshalUnmarshal проверяет корректность JSON-marshalling и unmarshalling для структуры Task.
// Проверяется, что сериализация содержит все поля и что после десериализации данные совпадают.
func TestTaskJSON_MarshalUnmarshal(t *testing.T) {
	id := uuid.New()
	now := time.Now().Truncate(time.Second)
	// Оригинальная задача для теста
	taskOrig := &Task{
		ID:         id,
		Status:     StatusCompleted,
		CreatedAt:  now,
		StartedAt:  &now,
		FinishedAt: &now,
		Result:     "ok",
		Error:      "",
	}

	// Сериализуем в JSON
	data, err := json.Marshal(taskOrig)
	if err != nil {
		t.Fatalf("не удалось выполнить Marshal: %v", err)
	}

	// Десериализуем обратно
	var taskClone Task
	if err := json.Unmarshal(data, &taskClone); err != nil {
		t.Fatalf("не удалось выполнить Unmarshal: %v", err)
	}

	// Сравниваем основные поля
	if taskClone.ID != taskOrig.ID {
		t.Errorf("ID mismatch: got %v, want %v", taskClone.ID, taskOrig.ID)
	}
	if taskClone.Status != taskOrig.Status {
		t.Errorf("Status mismatch: got %v, want %v", taskClone.Status, taskOrig.Status)
	}
	if !taskClone.CreatedAt.Equal(taskOrig.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", taskClone.CreatedAt, taskOrig.CreatedAt)
	}
	if taskClone.Result != taskOrig.Result {
		t.Errorf("Result mismatch: got %v, want %v", taskClone.Result, taskOrig.Result)
	}
}
