package shared

import (
	"encoding/json"
	"net/http"
)

type ErrorCode string

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

const (
	ErrBadRequest     ErrorCode = "bad_request"
	ErrInvalidInput   ErrorCode = "invalid_input"
	ErrUnauthorized   ErrorCode = "unauthorized"
	ErrInternalServer ErrorCode = "internal_server_error"
	ErrNotFound       ErrorCode = "not_found"
)

type APIResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message,omitempty"`
	Data    any               `json:"data,omitempty"`
	Error   ErrorCode         `json:"error,omitempty"`
	Errors  []ValidationError `json:"errors,omitempty"`
	Success bool              `json:"success"`
}

func Success(w http.ResponseWriter, status int, data any, message string) {
	send(w, status, APIResponse{
		Status:  status,
		Data:    data,
		Message: message,
		Success: true,
	})
}

func BadRequest(w http.ResponseWriter, err ErrorCode, validationErrors []ValidationError) {
	send(w, http.StatusBadRequest, APIResponse{
		Status:  http.StatusBadRequest,
		Error:   err,
		Errors:  validationErrors,
		Success: false,
	})
}

// Never expose internal database/message queues/system errors to the client. EVER.
func InternalError(w http.ResponseWriter) {
	send(w, http.StatusInternalServerError, APIResponse{
		Status:  http.StatusInternalServerError,
		Error:   "internal server error",
		Success: false,
	})
}

func send(w http.ResponseWriter, status int, res APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
