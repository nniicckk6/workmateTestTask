package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"workmateTestProject/internal/model"
)

// TestInMemoryTaskStore_Crud проверяет базовые операции Create, Get, Delete и List в InMemoryTaskStore.
// Каждый шаг сопровождается тщательным описанием целей и ожидаемых результатов.
func TestInMemoryTaskStore_Crud(t *testing.T) {
	// Создаём новое хранилище
	s := NewInMemoryTaskStore()
	// Проверяем, что список задач пуст
	initial := s.List()
	if len(initial) != 0 {
		t.Fatalf("ожидалось пустое хранилище, получили %d элементов", len(initial))
	}

	// Создаём новую задачу и сохраняем
	id := uuid.New()
	task := &model.Task{ID: id, Status: model.StatusPending, CreatedAt: time.Now()}
	s.Create(task)

	// Проверяем, что Get возвращает сохранённую задачу
	task2, ok := s.Get(id)
	if !ok {
		t.Fatalf("задача с ID %v не найдена после Create", id)
	}
	if task2.ID != id {
		t.Errorf("ID mismatch: got %v, want %v", task2.ID, id)
	}

	// Проверяем, что List возвращает один элемент
	all := s.List()
	if len(all) != 1 {
		t.Errorf("ожидалось 1 задачу в списке, получили %d", len(all))
	}

	// Удаляем задачу и проверяем, что она пропала
	s.Delete(id)
	_, ok = s.Get(id)
	if ok {
		t.Errorf("задача с ID %v всё ещё присутствует после Delete", id)
	}

	// Проверяем, что хранилище снова пустое
	if len(s.List()) != 0 {
		t.Errorf("ожидалось пустое хранилище после удаления, но есть элементы")
	}
}

// TestInMemoryTaskStore_Cancel проверяет метод Cancel: статус меняется на Canceled только для существующей задачи.
func TestInMemoryTaskStore_Cancel(t *testing.T) {
	s := NewInMemoryTaskStore()

	// Создаём задачу и добавляем в хранилище
	id := uuid.New()
	task := &model.Task{ID: id, Status: model.StatusPending, CreatedAt: time.Now()}
	s.Create(task)

	// Отменяем задачу, ожидаем true
	ok := s.Cancel(id)
	if !ok {
		t.Fatalf("Cancel вернул false для существующей задачи %v", id)
	}
	// Проверяем, что статус изменился
	task2, _ := s.Get(id)
	if task2.Status != model.StatusCanceled {
		t.Errorf("ожидалось StatusCanceled, получили %v", task2.Status)
	}

	// Повторный вызов Cancel для той же задачи тоже возвращает true
	if !s.Cancel(id) {
		t.Errorf("повторный Cancel вернул false, ожидался true")
	}

	// Cancel для несуществующего ID возвращает false
	if s.Cancel(uuid.New()) {
		t.Errorf("Cancel вернул true для несуществующего ID")
	}
}
