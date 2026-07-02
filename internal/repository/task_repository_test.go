package repository

import (
	"testing"
	"time"

	"github.com/NITHISH-2006/taskflow-go/internal/model"
)

func TestInMemoryTaskRepositoryCRUD(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	now := time.Now().UTC()
	task := model.Task{
		ID:          "task-1",
		Title:       "Test task",
		Description: "task description",
		Status:      model.StatusPending,
		Priority:    3,
		DueDate:     &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := repo.Create(task)
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != task.ID {
		t.Fatalf("expected id %q, got %q", task.ID, created.ID)
	}

	retrieved, err := repo.Get(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if retrieved.Title != task.Title {
		t.Fatalf("expected title %q, got %q", task.Title, retrieved.Title)
	}

	updated := task
	updated.Title = "Updated task"
	if _, err := repo.Update(updated); err != nil {
		t.Fatal(err)
	}

	found, err := repo.Get(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.Title != "Updated task" {
		t.Fatalf("expected updated title, got %q", found.Title)
	}

	if err := repo.Delete(task.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := repo.Get(task.ID); err != ErrTaskNotFound {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestInMemoryTaskRepositoryFiltering(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	past := time.Now().Add(-2 * time.Hour).UTC()
	future := time.Now().Add(2 * time.Hour).UTC()

	tasks := []model.Task{
		{ID: "task-1", Title: "Pending", Status: model.StatusPending, Priority: 2, DueDate: &past, CreatedAt: past, UpdatedAt: past},
		{ID: "task-2", Title: "Completed", Status: model.StatusCompleted, Priority: 5, DueDate: &future, CreatedAt: past, UpdatedAt: past},
	}

	for _, task := range tasks {
		if _, err := repo.Create(task); err != nil {
			t.Fatal(err)
		}
	}

	status := model.StatusPending
	filtered, err := repo.GetAll(model.TaskFilter{Status: &status})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].ID != "task-1" {
		t.Fatalf("expected task-1, got %v", filtered)
	}

	priority := 5
	filtered, err = repo.GetAll(model.TaskFilter{Priority: &priority})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].ID != "task-2" {
		t.Fatalf("expected task-2, got %v", filtered)
	}

	dueBefore := time.Now().UTC()
	filtered, err = repo.GetAll(model.TaskFilter{DueBefore: &dueBefore})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].ID != "task-1" {
		t.Fatalf("expected overdue task-1, got %v", filtered)
	}
}
