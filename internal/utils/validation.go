package utils

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrTitleRequired      = errors.New("title is required")
	ErrPriorityOutOfRange = errors.New("priority must be between 1 and 5")
	ErrStatusInvalid      = errors.New("status must be pending, in_progress, or completed")
	ErrDueDateInvalid     = errors.New("due_date must be a valid RFC3339 timestamp")
)

func ValidateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrTitleRequired
	}
	return nil
}

func ValidatePriority(priority int) error {
	if priority < 1 || priority > 5 {
		return ErrPriorityOutOfRange
	}
	return nil
}

func ValidateStatus(status string) error {
	if strings.TrimSpace(status) == "" {
		return nil
	}

	normalized := strings.ToLower(strings.TrimSpace(status))
	switch normalized {
	case "pending", "in_progress", "completed":
		return nil
	default:
		return ErrStatusInvalid
	}
}

func ValidateTimestamp(timestamp string) error {
	if strings.TrimSpace(timestamp) == "" {
		return nil
	}
	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		return ErrDueDateInvalid
	}
	return nil
}
