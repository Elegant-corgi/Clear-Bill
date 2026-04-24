package httpx

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func WriteData[T any](w http.ResponseWriter, status int, data T) {
	writeJSON(w, status, Response[T]{
		Success: true,
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, Response[any]{
		Success: false,
		Error:   message,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(payload)
}
