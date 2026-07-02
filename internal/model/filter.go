package model

import "time"

type TaskFilter struct {
	Status    *Status
	Priority  *int
	DueBefore *time.Time
}

func (filter TaskFilter) Matches(task Task) bool {
	if filter.Status != nil && task.Status != *filter.Status {
		return false
	}

	if filter.Priority != nil && task.Priority != *filter.Priority {
		return false
	}

	if filter.DueBefore != nil {
		if task.DueDate == nil {
			return false
		}
		if task.DueDate.After(filter.DueBefore.UTC()) || task.DueDate.Equal(filter.DueBefore.UTC()) {
			return false
		}
	}

	return true
}
