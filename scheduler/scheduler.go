package scheduler

import (
	"github.com/google/uuid"
	"github.com/tommsawyer/fetch/task"
)

// TasksStorage represents tasks storage.
type TasksStorage interface {
	Save(t *task.Task)
	Remove(id uuid.UUID)
	Get(id uuid.UUID) *task.Task
	All() []*task.Task
}

// Scheduler schedules tasks.
type Scheduler struct {
	store          TasksStorage
	scheduledTasks chan *task.Task
}

// New returns a new scheduler.
func New() *Scheduler {
	return &Scheduler{
		store:          NewInMemoryTasksStore(),
		scheduledTasks: make(chan *task.Task),
	}
}

// Run runs scheduler.
func (s *Scheduler) Run() {
	for t := range s.scheduledTasks {
		go t.Run()
	}
}

// Schedule creates task and schedules it.
func (s *Scheduler) Schedule(method, url, body string) *task.Task {
	t := task.New(method, url, body)
	s.store.Save(t)
	s.scheduledTasks <- t
	return t
}

// TaskByID returns task by id.
func (s *Scheduler) TaskByID(id uuid.UUID) *task.Task {
	return s.store.Get(id)
}

// Delete removes task.
// It will cancel task if task was in progress.
func (s *Scheduler) Delete(id uuid.UUID) {
	task := s.store.Get(id)
	if task != nil {
		task.Cancel()
	}
	s.store.Remove(id)
}

// Tasks returns all tasks.
func (s *Scheduler) Tasks() []*task.Task {
	return s.store.All()
}
