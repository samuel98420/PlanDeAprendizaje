package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"github.com/samuel98420/PlanDeAprendizaje/backend/internal/database"
	"github.com/samuel98420/PlanDeAprendizaje/backend/internal/logger"
)


func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string, details ...string) {
	response := map[string]interface{}{
		"error":   message,
		"details": details,
	}
	respondJSON(w, status, response)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.GetTasks()
	if err != nil {
		logger.Logger.WithField("error", err).Error("Error en GetTasks")
		respondError(w, http.StatusInternalServerError, "Error al obtener tareas", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logger.Logger.Warn("JSON inválido:", err)
		respondError(w, http.StatusBadRequest, "Formato de tarea inválido")
		return
	}

	if task.Title == "" {
		respondError(w, http.StatusBadRequest, "El título es obligatorio")
		return
	}

	id, err := database.CreateTask(task)
	if err != nil {
		logger.Logger.WithField("task", task).Error("Error al crear:", err)
		respondError(w, http.StatusInternalServerError, "Error al crear tarea", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "Tarea creada exitosamente",
	})
}

func HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetTaskByIDHandler(w, r, id)
	case http.MethodPut:
		UpdateTaskHandler(w, r, id)
	case http.MethodDelete:
		DeleteTaskHandler(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}

func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request, id int) {
	task, err := database.GetTaskByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Tarea no encontrada")
		return
	}
	respondJSON(w, http.StatusOK, task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	var task database.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondError(w, http.StatusBadRequest, "Formato inválido")
		return
	}

	if err := database.UpdateTask(id, task); err != nil {
		respondError(w, http.StatusInternalServerError, "Error al actualizar")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Tarea actualizada"})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	if err := database.DeleteTask(id); err != nil {
		respondError(w, http.StatusInternalServerError, "Error al eliminar")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Tarea eliminada"})
}