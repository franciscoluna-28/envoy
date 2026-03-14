package environments

import (
	"net/http"
	response "newserver/internal/shared"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	repo      Repository
	validator *validator.Validate
	key       []byte
}

func NewHandler(repo Repository, v *validator.Validate, key []byte) *Handler {
	return &Handler{
		repo:      repo,
		validator: v,
		key:       key,
	}
}

func (h *Handler) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	var req CreateEnvironmentRequest

	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	env := CreateProjectEnvironment(r.Context(), req, h.key, h.repo)

	response.WriteJSON(w, http.StatusCreated, env)
}
