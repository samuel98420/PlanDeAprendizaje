package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Definición de la estructura Task
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// initDB inicializa la base de datos y crea la tabla tasks
func initDB() {
	// Abre la conexión con la base de datos SQLite
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query para crear la tabla de tareas
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error al crear la tabla: %v", err)
	}

	fmt.Println("Base de datos inicializada correctamente")
}

// createTask agrega una nueva tarea a la base de datos
func createTask(title, description string, completed bool) (int, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer db.Close()

	// Inserta una nueva tarea en la base de datos
	result, err := db.Exec("INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)", title, description, completed)
	if err != nil {
		log.Printf("Error al insertar tarea: %v", err)
		return 0, err
	}

	// Obtiene el ID de la tarea recién creada
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error al obtener el ID de la nueva tarea: %v", err)
		return 0, err
	}

	return int(lastID), nil
}

// getTasks obtiene todas las tareas almacenadas en la base de datos
func getTasks() ([]Task, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	// Ejecuta la consulta para obtener todas las tareas
	rows, err := db.Query("SELECT id, title, description, completed FROM tasks")
	if err != nil {
		log.Printf("Error al ejecutar la consulta: %v", err)
		return nil, fmt.Errorf("error al obtener tareas")
	}
	defer rows.Close()

	var tasks []Task

	// Recorre los resultados y los almacena en la lista de tareas
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed); err != nil {
			log.Printf("Error al escanear fila: %v", err)
			return nil, fmt.Errorf("error al procesar los datos de tareas")
		}
		tasks = append(tasks, task)
	}

	// Verifica si ocurrió un error durante la iteración
	if err = rows.Err(); err != nil {
		log.Printf("Error al recorrer filas: %v", err)
		return nil, fmt.Errorf("error al leer datos")
	}

	return tasks, nil
}

// getTaskByID obtiene una tarea específica según su ID
func getTaskByID(id int) (Task, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
		return Task{}, err
	}
	defer db.Close()

	var task Task

	// Consulta para obtener la tarea con el ID
	err = db.QueryRow("SELECT id, title, description, completed FROM tasks WHERE id = ?", id).
		Scan(&task.ID, &task.Title, &task.Description, &task.Completed)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

// updateTask actualiza los datos de una tarea
func updateTask(id int, title, description string, completed bool) (int64, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer db.Close()

	// Ejecuta la actualización de la tarea
	result, err := db.Exec("UPDATE tasks SET title = ?, description = ?, completed = ? WHERE id = ?", title, description, completed, id)
	if err != nil {
		log.Printf("Error al actualizar tarea con ID %d: %v", id, err)
		return 0, err
	}

	// Obtiene el número de filas afectadas por la actualización
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// deleteTask elimina una tarea de la base de datos según su ID
func deleteTask(id int) (int64, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer db.Close()

	log.Printf("🔍 Ejecutando consulta DELETE para ID: %d", id)

	// Ejecuta la eliminación de la tarea
	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Printf("Error al eliminar la tarea con ID %d: %v", id, err)
		return 0, err
	}

	// Obtiene el número de filas afectadas por la eliminación
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
