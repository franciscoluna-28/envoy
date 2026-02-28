package infra

import (
	"encoding/json"
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
	return json.NewDecoder(r.Body).Decode(v)
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
