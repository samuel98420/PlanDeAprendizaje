package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"


)

var db *sql.DB
var log = logrus.WithField("component", "database")

type Task struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}

func InitDB() error {
	dbPath := "/Users/macbookair2020/documents/PlanDeAprendizaje-main/backend/internal/database/tasks.db"
	log.Infof("📂 Ruta de la DB: %s", dbPath)

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Errorf("❌ Error al abrir DB: %v", err)
		return err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0
	);`
	if _, err = db.Exec(createTableQuery); err != nil {
		log.Errorf("💥 Error al crear tabla: %v", err)
		return fmt.Errorf("error al crear la tabla: %v", err)
	}

	log.Info("✅ Base de datos lista")
	return nil
}

func CreateTask(task Task) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)",
		task.Title, task.Description, task.Completed,
	)
	if err != nil {
		return 0, fmt.Errorf("error al crear tarea: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener ID: %v", err)
	}

	return id, nil
}

// GetTasks obtiene todas las tareas de la DB
func GetTasks() ([]Task, error) {
	rows, err := db.Query("SELECT id, title, description, completed FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("error al consultar tareas: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed); err != nil {
			return nil, fmt.Errorf("error al leer tarea: %v", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error final al leer tareas: %v", err)
	}

	return tasks, nil
}

// GetTaskByID obtiene una tarea por su ID
func GetTaskByID(id int) (Task, error) {
	var task Task
	err := db.QueryRow(
		"SELECT id, title, description, completed FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Completed)

	if err != nil {
		return Task{}, fmt.Errorf("error al obtener tarea %d: %v", id, err)
	}

	return task, nil
}

// UpdateTask con transacción
func UpdateTask(id int, task Task) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer tx.Rollback() 

	_, err = tx.Exec(
		"UPDATE tasks SET title = ?, description = ?, completed = ? WHERE id = ?",
		task.Title, task.Description, task.Completed, id,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %v", err)
	}
	return nil
}

// DeleteTask elimina una tarea por ID
func DeleteTask(id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error al eliminar tarea %d: %v", id, err)
	}
	return nil
}

// CloseDB cierra la conexión con la DB (llamar al final del programa)
func CloseDB() {
	db.Close()
	log.Println("🔌 Conexión a DB cerrada")
}


