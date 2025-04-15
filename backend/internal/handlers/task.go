package handlers

import (
	"encoding/json"
	"github.com/tuusuario/todo-api/backend/internal/database"
	"github.com/tuusuario/todo-api/backend/internal/models"
    "github.com/tuusuario/todo-api/backend/internal/errors"
	"net/http"
	"strconv"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.GetTasks()
	if err != nil {
		errors.RespondWithError(w, http.StatusInternalServerError, 
			"Error al obtener tareas", 
			err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Formato de tarea inválido", http.StatusBadRequest)
		return
	}

	id, err := database.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Tarea creada exitosamente",
	})
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID debe ser un número", http.StatusBadRequest)
		return
	}

	task, err := database.GetTaskByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID debe ser un número", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Formato de tarea inválido", http.StatusBadRequest)
		return
	}

	if err := database.UpdateTask(id, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Tarea actualizada exitosamente",
	})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID debe ser un número", http.StatusBadRequest)
		return
	}

	if err := database.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Tarea eliminada exitosamente",
	})
}
