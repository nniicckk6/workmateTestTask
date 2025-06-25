package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus представляет статус задачи
// Допустимые значения: Pending, InProgress, Completed, Failed, Canceled
type TaskStatus string

const (
	StatusPending    TaskStatus = "Pending"
	StatusInProgress TaskStatus = "InProgress"
	StatusCompleted  TaskStatus = "Completed"
	StatusFailed     TaskStatus = "Failed"
	StatusCanceled   TaskStatus = "Canceled"
)

// Task описывает I/O-bound задачу
type Task struct {
	ID         uuid.UUID  `json:"id"`
	Status     TaskStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Result     string     `json:"result,omitempty"`
	Error      string     `json:"error,omitempty"`
}
