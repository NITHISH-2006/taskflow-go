package service

import (
	"strings"
	"time"

	"github.com/NITHISH-2006/taskflow-go/internal/model"
	"github.com/NITHISH-2006/taskflow-go/internal/repository"
	"github.com/NITHISH-2006/taskflow-go/internal/utils"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(req model.CreateTaskRequest) (model.Task, error) {
	if err := utils.ValidateTitle(req.Title); err != nil {
		return model.Task{}, err
	}

	if err := utils.ValidatePriority(req.Priority); err != nil {
		return model.Task{}, err
	}

	if err := utils.ValidateStatus(req.Status); err != nil {
		return model.Task{}, err
	}

	if err := utils.ValidateTimestamp(req.DueDate); err != nil {
		return model.Task{}, err
	}

	dueDate, err := parseDueDate(req.DueDate)
	if err != nil {
		return model.Task{}, err
	}

	status := normalizeStatus(req.Status)
	now := time.Now().UTC()
	task := model.Task{
		ID:          utils.NewID(),
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      status,
		Priority:    req.Priority,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Create(task)
}

func (s *TaskService) GetTask(id string) (model.Task, error) {
	return s.repo.Get(id)
}

func (s *TaskService) ListTasks(filter model.TaskFilter) ([]model.Task, error) {
	return s.repo.GetAll(filter)
}

func (s *TaskService) UpdateTask(id string, req model.UpdateTaskRequest) (model.Task, error) {
	task, err := s.repo.Get(id)
	if err != nil {
		return model.Task{}, err
	}

	if req.Title != nil {
		if err := utils.ValidateTitle(*req.Title); err != nil {
			return model.Task{}, err
		}
		task.Title = strings.TrimSpace(*req.Title)
	}

	if req.Description != nil {
		task.Description = strings.TrimSpace(*req.Description)
	}

	if req.Status != nil {
		if err := utils.ValidateStatus(*req.Status); err != nil {
			return model.Task{}, err
		}
		task.Status = normalizeStatus(*req.Status)
	}

	if req.Priority != nil {
		if err := utils.ValidatePriority(*req.Priority); err != nil {
			return model.Task{}, err
		}
		task.Priority = *req.Priority
	}

	if req.DueDate != nil {
		trimmed := strings.TrimSpace(*req.DueDate)
		if trimmed == "" {
			task.DueDate = nil
		} else {
			if err := utils.ValidateTimestamp(trimmed); err != nil {
				return model.Task{}, err
			}
			parsed, err := parseDueDate(trimmed)
			if err != nil {
				return model.Task{}, err
			}
			task.DueDate = parsed
		}
	}

	task.UpdatedAt = time.Now().UTC()
	return s.repo.Update(task)
}

func (s *TaskService) DeleteTask(id string) error {
	return s.repo.Delete(id)
}

func normalizeStatus(input string) model.Status {
	if strings.TrimSpace(input) == "" {
		return model.StatusPending
	}

	normalized := strings.ToLower(strings.TrimSpace(input))
	switch normalized {
	case string(model.StatusPending):
		return model.StatusPending
	case string(model.StatusInProgress):
		return model.StatusInProgress
	case string(model.StatusCompleted):
		return model.StatusCompleted
	default:
		return model.StatusPending
	}
}

func parseDueDate(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	parsed, err := utils.ParseTime(value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
