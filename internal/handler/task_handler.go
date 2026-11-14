package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/anfeespi/devboard/internal/domain"
	"github.com/anfeespi/devboard/internal/usecase"
	"github.com/anfeespi/devboard/internal/validator"
)

// TaskHandler struct
type TaskHandler struct {
	usecase   *usecase.TaskUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

type updateStatusRequest struct {
	Status domain.TaskStatus `json:"status" validate:"required"`
}

type assignTaskRequest struct {
	AssigneeID string `json:"assignee_id" validate:"required,uuid4"`
}

// NewTaskHandler constructor
func NewTaskHandler(uc *usecase.TaskUseCase, v *validator.Validator, l *slog.Logger) *TaskHandler {
	return &TaskHandler{
		usecase:   uc,
		validator: v,
		logger:    l,
	}
}

// createTaskRequest struct
type createTaskRequest struct {
	ProjectID   string `json:"project_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"max=2000"`
}

// taskResponse struct
type taskResponse struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	Status      domain.TaskStatus `json:"status"`
	AssigneeID  *string           `json:"assignee_id,omitempty"`
	CreatedBy   string            `json:"created_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// toTaskResponse map the entity to the DTO
func toTaskResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		AssigneeID:  t.AssigneeID,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// Create create a new task endpoint
//
// @Summary Create task
// @Description Register a new task in the system
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body createTaskRequest true "Task data"
// @Success 201 {object} taskResponse
// @Failure 400 {object} ProblemDetails
// @Router /tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"error writing validation errors",
				"error", err,
				"method", r.Method,
				"path", r.URL.Path,
			)
		}
		return
	}

	task, err := h.usecase.CreateTask(r.Context(), req.ProjectID, req.Title, req.Description, "system")
	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error writing validation errors",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}
}

// Get get a task by id
//
// @Summary Get task by id
// @Description Get a task using its id
// @Tags tasks
// @Produce json
// @Param id path string true "task id"
// @Success 200 {object} taskResponse
// @Failure 404 {object} ProblemDetails
// @Router /tasks/{id} [get]
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	task, err := h.usecase.GetTask(r.Context(), id)
	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(r.Context(), "error getting task", "error", err, "method", r.Method, "path", r.URL.Path)
	}
}

// ListByProject get all the tasks of one project
//
// @Summary Get tasks by project
// @Description Get all the tasks based on its project
// @Tags tasks
// @Produce json
// @Param id path string true "project id"
// @Success 200 {object} PaginatedResponse[taskResponse]
// @Failure 404 {object} ProblemDetails
// @Router /projects/{id}/tasks [get]
func (h *TaskHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	tasks, err := h.usecase.ListProjectTasks(r.Context(), projectID)
	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	response := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		response[i] = toTaskResponse(t)
	}

	if err := RespondJSON(w, http.StatusOK, PaginatedResponse[taskResponse]{
		Data:    response,
		HasMore: false,
	}); err != nil {
		h.logger.ErrorContext(r.Context(), "error getting the tasks by project", "error", err, "method", r.Method, "path", r.URL.Path)
	}
}

// UpdateStatus update task status
//
// @Summary update task status
// @Description Move the task status to a new one
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task id"
// @Param request body updateStatusRequest true "new status"
// @Success 200 {object} taskResponse
// @Failure 400 {object} ProblemDetails
// @Router /tasks/{id}/status [put]
func (h *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req updateStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"error writing validation errors",
				"error", err,
				"method", r.Method,
				"path", r.URL.Path,
			)
		}
		return
	}

	task, err := h.usecase.UpdateTaskStatus(r.Context(), id, req.Status)
	if err != nil {
		RespondError(w, r, h.logger, err)
	}

	if err := RespondJSON(w, http.StatusOK, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(r.Context(), "error updating task", "error", err, "method", r.Method, "path", r.URL.Path)
	}
}

// Assign assign task to user
//
// @Summary assign task to user
// @Description Assign a task to a user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task id"
// @Param request body assignTaskRequest true "assignee id"
// @Success 200 {object} taskResponse
// @Failure 400 {object} ProblemDetails
// @Router /tasks/{id}/assignee [put]
func (h *TaskHandler) Assign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req assignTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"error writing validation errors",
				"error", err,
				"method", r.Method,
				"path", r.URL.Path,
			)
		}
		return
	}

	task, err := h.usecase.AssignTask(r.Context(), id, req.AssigneeID)
	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(r.Context(), "error assigning task", "error", err, "method", r.Method, "path", r.URL.Path)
	}
}
