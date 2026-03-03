package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func Action(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
}

func Decode(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is required")
	}
	if r.ContentLength == 0 {
		return fmt.Errorf("request body cannot be empty")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
