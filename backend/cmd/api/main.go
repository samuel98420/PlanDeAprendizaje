package main

import (
	"context"
	"net/http"
	"github.com/tuusuario/todo-api/backend/internal/database"
	"github.com/tuusuario/todo-api/backend/internal/handlers"
	"github.com/tuusuario/todo-api/backend/internal/logger"
)

func main() {
	logger.Init() // Inicializa logger
	
	if err := database.InitDB(); err != nil {
		logger.Logger.Fatalf("Failed to init DB: %v", err)
	}
	defer database.CloseDB()

	// Middleware para inyectar logger
	withLogger := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger.LogRequest(r)
			ctx := context.WithValue(r.Context(), "logger", logger.Logger)
			next(w, r.WithContext(ctx))
		}
	}

	// Rutas
	http.HandleFunc("/tasks", withLogger(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTasks(w, r)
		// ... otros métodos
		}
	}))
	
	logger.Logger.Info("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}