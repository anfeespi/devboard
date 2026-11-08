// Package domain domain
package domain

import (
	"encoding/json"
	"fmt"
)

// TaskStatus status of the task
type TaskStatus string

// Enum
const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// IsValid validate a task string
func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

// UnmarshalJSON deserialize a json
func (s *TaskStatus) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	status := TaskStatus(raw)
	if !status.IsValid() {
		return fmt.Errorf("task status invalid: %q", raw)
	}

	*s = status
	return nil
}
