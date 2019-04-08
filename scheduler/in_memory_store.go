package scheduler

import (
	"sync"

	"github.com/google/uuid"
	"github.com/tommsawyer/fetch/task"
)

// InMemoryTasksStore stores tasks in memory.
type InMemoryTasksStore struct {
	mu    sync.RWMutex
	tasks map[string]*task.Task
}

// NewInMemoryTasksStore creates a new memory tasks store.
func NewInMemoryTasksStore() *InMemoryTasksStore {
	return &InMemoryTasksStore{
		tasks: make(map[string]*task.Task),
	}
}

// Save saves task.
func (s *InMemoryTasksStore) Save(task *task.Task) {
	s.mu.Lock()
	s.tasks[task.ID.String()] = task
	s.mu.Unlock()
}

// Remove removes task.
func (s *InMemoryTasksStore) Remove(id uuid.UUID) {
	s.mu.Lock()
	delete(s.tasks, id.String())
	s.mu.Unlock()
}

// Get returns task by id.
func (s *InMemoryTasksStore) Get(id uuid.UUID) *task.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[id.String()]
}

// All returns all tasks.
func (s *InMemoryTasksStore) All() []*task.Task {
	s.mu.RLock()
	tasks := make([]*task.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	s.mu.RUnlock()
	return tasks
}
