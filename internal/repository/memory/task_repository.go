package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/anfeespi/devboard/internal/domain"
)

// TaskRepository struct
type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

// NewTaskRepository constructor
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

// Create method to create a task
func (r *TaskRepository) Create(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = fmt.Sprintf("task_%d", len(r.tasks)+1)
	r.tasks[task.ID] = task

	return nil
}

// Update method to update a task
func (r *TaskRepository) Update(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task

	return nil
}

// GetByID get a task by its ID
func (r *TaskRepository) GetByID(_ context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, fmt.Errorf("tarea %s: %w", id, domain.ErrNotFound)
	}

	return task, nil
}

// ListByProject list tasks by its project
func (r *TaskRepository) ListByProject(_ context.Context, projectID string) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasksByProject := make([]*domain.Task, 0)
	for _, v := range r.tasks {
		if v.ProjectID == projectID {
			tasksByProject = append(tasksByProject, v)
		}
	}

	return tasksByProject, nil
}
