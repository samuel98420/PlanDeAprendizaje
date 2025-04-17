package errors

import (
	"encoding/json"
	"net/http"
	"runtime"
)

// ApiError representa un error estándar de la API
type ApiError struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	Detail     string `json:"detail,omitempty"`
	Location   string `json:"location,omitempty"` 
}

// RespondWithError envía errores en formato JSON
func RespondWithError(w http.ResponseWriter, status int, message string, detail string) {
	_, file, line, _ := runtime.Caller(1) 
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiError{
		Status:   status,
		Message:  message,
		Detail:   detail,
		Location: file + ":" + string(line),
	})
}