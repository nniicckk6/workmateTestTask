package service

import (
	"fmt"
	"math/rand"
	"time"

	"workmateTestProject/internal/model"
)

// SimulateWorkFunc указывает на функцию-симулятор, может быть переопределена в тестах
var SimulateWorkFunc = simulateWork

// simulateWork симулирует I/O-bound работу, возвращая результат или ошибку
func simulateWork() (string, error) {
	// Ждем случайное время от 1 до 5 минут
	dur := time.Duration(rand.Intn(5)+1) * time.Minute
	time.Sleep(dur)
	// Возвращаем результат
	return fmt.Sprintf("Обработано за %s", dur), nil
}

// StartProcessing запускает обработку задачи в отдельной горутине
func StartProcessing(task *model.Task) {
	go func() {
		task.Status = model.StatusInProgress
		now := time.Now()
		task.StartedAt = &now

		result, err := SimulateWorkFunc()
		finish := time.Now()
		task.FinishedAt = &finish
		if err != nil {
			task.Status = model.StatusFailed
			task.Error = err.Error()
		} else {
			task.Status = model.StatusCompleted
			task.Result = result
		}
	}()
}
