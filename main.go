package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// initLogging configura el sistema de registro de eventos para guardar logs en un archivo
func initLogging() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf(" No se pudo abrir el archivo de log: %v", err)
	}

	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Servidor iniciado y registro de eventos activado")
}

func main() {
	initDB()
	initLogging()

	// Definición de rutas para el manejo de tareas
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/", handleTaskByID)

	fmt.Println("API corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleTasks maneja solicitudes HTTP relacionadas con la lista de tareas
func handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTasksHandler(w, r)
	case "POST":
		createTaskHandler(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// handleTaskByID maneja operaciones CRUD sobre tareas individuales
func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID no encontrado en la URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		log.Printf("Error al convertir ID: %s, error: %v", idStr, err)
		return
	}

	switch r.Method {
	case "GET":
		getTaskByIDHandler(w, r, id)
	case "PUT":
		updateTaskHandler(w, r, id)
	case "DELETE":
		deleteTaskHandler(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// createTaskHandler maneja la creación de una nueva tarea
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error": "Formato JSON inválido"}`, http.StatusBadRequest)
		log.Printf("Error al decodificar JSON: %v", err)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "El campo 'title' es obligatorio"}`, http.StatusBadRequest)
		return
	}

	id, err := createTask(task.Title, task.Description, task.Completed)
	if err != nil {
		http.Error(w, `{"error": "No se pudo crear la tarea"}`, http.StatusInternalServerError)
		log.Printf("Error al insertar tarea en la base de datos: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tarea creada con éxito",
		"id":      id,
	})
}

// getTasksHandler obtiene todas las tareas
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := getTasks()
	if err != nil {
		http.Error(w, "No se pudieron obtener las tareas", http.StatusInternalServerError)
		log.Printf("Error al obtener tareas: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// getTaskByIDHandler obtiene una tarea específica por ID
func getTaskByIDHandler(w http.ResponseWriter, r *http.Request, id int) {
	task, err := getTaskByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "No se encontró la tarea con ID %d"}`, id), http.StatusNotFound)
		log.Printf("Error al obtener tarea con ID %d: %v", id, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// updateTaskHandler maneja la actualización de una tarea existente
func updateTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error": "Formato JSON inválido"}`, http.StatusBadRequest)
		log.Printf("Error al decodificar JSON: %v", err)
		return
	}
	rowsAffected, err := updateTask(id, task.Title, task.Description, task.Completed)

	if err != nil {
		http.Error(w, `{"error": "Error al actualizar la tarea"}`, http.StatusInternalServerError)
		log.Printf("Error al actualizar tarea con ID %d: %v", id, err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error": "No se encontró la tarea o no hubo cambios"}`, http.StatusNotFound)
		log.Printf("No se encontró la tarea con ID %d o no hubo cambios", id)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Tarea actualizada correctamente"})
}

// deleteTaskHandler maneja la eliminación de una tarea
func deleteTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	rowsAffected, err := deleteTask(id)
	if err != nil {
		http.Error(w, `{"error": "Error al intentar eliminar la tarea"}`, http.StatusInternalServerError)
		log.Printf("Error al eliminar la tarea con ID %d: %v", id, err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error": "No se encontró la tarea para eliminar"}`, http.StatusNotFound)
		log.Printf("Intento de eliminar tarea inexistente con ID %d", id)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Tarea eliminada correctamente"})
}
