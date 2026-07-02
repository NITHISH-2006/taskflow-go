package repository

import (
	"errors"
	"sync"

	"github.com/NITHISH-2006/taskflow-go/internal/model"
)

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
)

type TaskRepository interface {
	Create(task model.Task) (model.Task, error)
	Get(id string) (model.Task, error)
	GetAll(filter model.TaskFilter) ([]model.Task, error)
	Update(task model.Task) (model.Task, error)
	Delete(id string) error
}

type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]model.Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]model.Task),
	}
}

func (r *InMemoryTaskRepository) Create(task model.Task) (model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return model.Task{}, ErrTaskAlreadyExists
	}

	r.tasks[task.ID] = task
	return task, nil
}

func (r *InMemoryTaskRepository) Get(id string) (model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return model.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (r *InMemoryTaskRepository) GetAll(filter model.TaskFilter) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]model.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		if filter.Matches(task) {
			list = append(list, task)
		}
	}

	return list, nil
}

func (r *InMemoryTaskRepository) Update(task model.Task) (model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return model.Task{}, ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return task, nil
}

func (r *InMemoryTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}
