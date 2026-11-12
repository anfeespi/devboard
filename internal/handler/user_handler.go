package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/anfeespi/devboard/internal/domain"
	"github.com/anfeespi/devboard/internal/usecase"
	"github.com/anfeespi/devboard/internal/validator"
)

// UserHandler translate http to the use case
type UserHandler struct {
	usecase   *usecase.UserUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

// NewUserHandler create a new user handler
func NewUserHandler(uc *usecase.UserUseCase, v *validator.Validator, l *slog.Logger) *UserHandler {
	return &UserHandler{
		usecase:   uc,
		validator: v,
		logger:    l,
	}
}

type createUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Create create a new user endpoint
//
// @Summary Create user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param request body createUserRequest true "User data"
// @Success 201 {object} userResponse
// @Failure 400 {object} ProblemDetails
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(r.Context(), "error writing validation response", "error", err, "method", r.Method, "path", r.URL.Path)
		}
		return
	}

	// temporal
	passwordHash := "TEMPORAL_" + req.Password

	user, err := h.usecase.CreateUser(r.Context(), req.Email, req.Name, passwordHash)

	if err != nil {
		RespondError(w, r, h.logger, err)
	}

	if err := RespondJSON(w, http.StatusCreated, userResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}); err != nil {
		h.logger.ErrorContext(r.Context(), "error writing json response", "error", err, "method", r.Method, "path", r.URL.Path)
	}
}

// Get user endpoint
//
// @Summary Obtain user
// @Description Get a user by its id
// @Tags users
// @Produce json
// @Param id path string true "ID of the user"
// @Success 200 {object} userResponse
// @Failure 404 {object} ProblemDetails
// @Router /users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.usecase.GetUser(r.Context(), id)
	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, userResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}); err != nil {
		h.logger.ErrorContext(r.Context(), "error getting user", "error", err, "method", r.Method, "path", r.URL.Path)
	}
}
