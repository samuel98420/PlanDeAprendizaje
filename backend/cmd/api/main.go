package main

import (
	"net/http"
	"time"  
	"github.com/samuel98420/PlanDeAprendizaje/backend/internal/database"
	"github.com/samuel98420/PlanDeAprendizaje/backend/internal/handlers"
	"github.com/samuel98420/PlanDeAprendizaje/backend/internal/logger"
)

func main() {
	logger.Init()

	if err := database.InitDB(); err != nil {
		logger.Logger.Fatal("Fallo crítico al iniciar DB:", err)
	}
	defer database.CloseDB()

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTasks(w, r)
		case http.MethodPost:
			handlers.CreateTask(w, r)
		}
	})

	http.HandleFunc("/tasks/", handlers.HandleTaskByID)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Logger.Info("🚀 Servidor iniciado en :8080")
	if err := server.ListenAndServe(); err != nil {
		logger.Logger.Fatal("Error en servidor:", err)
	}
}