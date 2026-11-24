package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/anfeespi/devboard/internal/domain"
	postgres "github.com/anfeespi/devboard/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TaskRepository implements use case of tasks
type TaskRepository struct {
	pool    *pgxpool.Pool
	queries *postgres.Queries
}

// NewTaskRepository constructor
func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool:    pool,
		queries: postgres.New(pool),
	}
}

// Create create a new task
func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	var projectID pgtype.UUID
	if err := projectID.Scan(task.ProjectID); err != nil {
		return fmt.Errorf("project_id invalid: %w", domain.ErrInvalidInput)
	}

	var createdBy pgtype.UUID
	if err := createdBy.Scan(task.CreatedBy); err != nil {
		return fmt.Errorf("created_by invalid: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.CreateTask(ctx, postgres.CreateTaskParams{
		ProjectID:   projectID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		CreatedBy:   createdBy,
	})
	if err != nil {
		return fmt.Errorf("inserting task in db: %w", err)
	}

	task.ID = result.ID.String()
	task.CreatedAt = result.CreatedAt.Time
	task.UpdatedAt = result.UpdatedAt.Time

	return nil
}

// Update update an existing task
func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	var id pgtype.UUID
	if err := id.Scan(task.ID); err != nil {
		return fmt.Errorf("id invalid: %w", domain.ErrInvalidInput)
	}

	var assigneeID pgtype.UUID
	if task.AssigneeID != nil {
		if err := assigneeID.Scan(*task.AssigneeID); err != nil {
			return fmt.Errorf("assignee_id invalid: %w", domain.ErrInvalidInput)
		}
	}

	_, err := r.queries.UpdateTask(ctx, postgres.UpdateTaskParams{
		ID:          id,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		AssigneeID:  assigneeID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("task %s: %w", task.ID, domain.ErrNotFound)
		}

		return fmt.Errorf("updating task %s: %w", task.ID, err)
	}

	return nil
}

// GetByID obtain a task by its id
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id %s: %w", id, domain.ErrNotFound)
	}

	result, err := r.queries.GetTaskByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task %s: %w", id, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("consulting task %s: %w", id, err)
	}

	return toDomainTask(result), nil
}

// ListByProject get tasks by project
func (r *TaskRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.Task, error) {
	var pgProjectID pgtype.UUID
	if err := pgProjectID.Scan(projectID); err != nil {
		return nil, fmt.Errorf("invalid project_id %s: %w", projectID, domain.ErrNotFound)
	}

	results, err := r.queries.ListTasksByProject(ctx, pgProjectID)
	if err != nil {
		return nil, fmt.Errorf("listing tasks by project %s: %w", projectID, err)
	}

	tasks := make([]*domain.Task, len(results))
	for i, t := range results {
		tasks[i] = toDomainTask(t)
	}

	return tasks, nil
}

// toDomainTask map the model sqlc to domain
func toDomainTask(t postgres.Task) *domain.Task {
	task := &domain.Task{
		ID:          t.ID.String(),
		ProjectID:   t.ProjectID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      domain.TaskStatus(t.Status),
		CreatedBy:   t.CreatedBy.String(),
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}

	if t.AssigneeID.Valid {
		id := t.AssigneeID.String()
		task.AssigneeID = &id
	}

	return task
}
