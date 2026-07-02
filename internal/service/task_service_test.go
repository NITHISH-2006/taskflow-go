package service

import (
	"testing"
	"time"

	"github.com/NITHISH-2006/taskflow-go/internal/model"
	"github.com/NITHISH-2006/taskflow-go/internal/repository"
)

func TestTaskService_CreateGetUpdateDelete(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	svc := NewTaskService(repo)

	dueDate := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	created, err := svc.CreateTask(model.CreateTaskRequest{
		Title:       "Learn Go",
		Description: "Complete the TaskFlow API",
		Status:      "in_progress",
		Priority:    3,
		DueDate:     dueDate,
	})
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	if created.ID == "" {
		t.Fatal("expected created task to have an ID")
	}

	if created.Status != model.StatusInProgress {
		t.Fatalf("expected status in_progress, got %q", created.Status)
	}

	fetched, err := svc.GetTask(created.ID)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}

	if fetched.Title != created.Title {
		t.Fatalf("expected title %q, got %q", created.Title, fetched.Title)
	}

	updated, err := svc.UpdateTask(created.ID, model.UpdateTaskRequest{
		Title:    ptrString("Learn Go and Test"),
		Priority: ptrInt(4),
	})
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	if updated.Title != "Learn Go and Test" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	if updated.Priority != 4 {
		t.Fatalf("expected updated priority 4, got %d", updated.Priority)
	}

	if err := svc.DeleteTask(created.ID); err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	if _, err := svc.GetTask(created.ID); err == nil {
		t.Fatal("expected GetTask to return error after delete")
	}
}

func TestTaskService_ListTasks_WithFilter(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	svc := NewTaskService(repo)

	_, err := svc.CreateTask(model.CreateTaskRequest{
		Title:    "Task One",
		Priority: 1,
	})
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	_, err = svc.CreateTask(model.CreateTaskRequest{
		Title:    "Task Two",
		Status:   "completed",
		Priority: 5,
	})
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	status := model.StatusCompleted
	filtered, err := svc.ListTasks(model.TaskFilter{Status: &status})
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}

	if len(filtered) != 1 || filtered[0].Status != model.StatusCompleted {
		t.Fatalf("expected one completed task, got %d", len(filtered))
	}
}

func ptrString(value string) *string {
	return &value
}

func ptrInt(value int) *int {
	return &value
}
