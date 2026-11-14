package domain

import (
	"errors"
	"time"
)

// Task represent a task in the app
type Task struct {
	ID          string
	ProjectID   string
	Title       string
	Description string
	Status      TaskStatus
	AssigneeID  *string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTask creates a new task
func NewTask(projectID, title, description, createdBy string) *Task {
	now := time.Now()

	return &Task{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      TaskStatusTodo,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

var (
	// ErrTaskAlreadyDone task already done
	ErrTaskAlreadyDone = errors.New("the task was already completed")
	//ErrTaskNotAssignable task not assignable to done
	ErrTaskNotAssignable = errors.New("task cannot be assigned in this status")
)

// AssignTo assign the task to a user
func (t *Task) AssignTo(userID string) error {
	if t.Status == TaskStatusDone {
		return ErrTaskNotAssignable
	}

	t.AssigneeID = &userID
	t.UpdatedAt = time.Now()
	return nil
}

// MarkAsDone mark the task as completed
func (t *Task) MarkAsDone() error {
	if t.Status == TaskStatusDone {
		return ErrTaskAlreadyDone
	}

	t.Status = TaskStatusDone
	t.UpdatedAt = time.Now()
	return nil
}

// MoveToStatus move the task to a determined status
func (t *Task) MoveToStatus(status TaskStatus) error {
	if !status.IsValid() {
		return errors.New("new status invalid")
	}

	t.Status = status
	t.UpdatedAt = time.Now()
	return nil
}
