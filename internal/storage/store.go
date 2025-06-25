package storage

import (
	"sync"

	"github.com/google/uuid"
	"workmateTestProject/internal/model"
)

// TaskStore определяет интерфейс потокобезопасного хранилища задач
type TaskStore interface {
	// Create добавляет новую задачу в хранилище
	Create(task *model.Task)
	// Get возвращает задачу по ID и флаг наличия
	Get(id uuid.UUID) (*model.Task, bool)
	// Delete удаляет задачу по ID
	Delete(id uuid.UUID)
	// List возвращает все задачи
	List() []*model.Task
	// Cancel устанавливает статус задачи Canceled, возвращает true если задача найдена
	Cancel(id uuid.UUID) bool
}

// InMemoryTaskStore - реализация TaskStore в памяти
type InMemoryTaskStore struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*model.Task
}

// NewInMemoryTaskStore создаёт новый InMemoryTaskStore
func NewInMemoryTaskStore() TaskStore {
	return &InMemoryTaskStore{
		tasks: make(map[uuid.UUID]*model.Task),
	}
}

// Create добавляет новую задачу в хранилище
func (s *InMemoryTaskStore) Create(task *model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

// Get возвращает задачу по ID
func (s *InMemoryTaskStore) Get(id uuid.UUID) (*model.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

// Delete удаляет задачу по ID
func (s *InMemoryTaskStore) Delete(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
}

// List возвращает все задачи
func (s *InMemoryTaskStore) List() []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		list = append(list, task)
	}
	return list
}

// Cancel устанавливает статус задачи Canceled
func (s *InMemoryTaskStore) Cancel(id uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return false
	}
	task.Status = model.StatusCanceled
	return true
}
