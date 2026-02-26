package repos

import (
	"database/sql"
	"log"
	"tasks-api/models"
)

type DB struct {
	db *sql.DB
}

func New(db *sql.DB) *DB {
	return &DB{db: db}
}

// init schema ensure table exists
func (r *DB) InitSchema() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            completed BOOLEAN DEFAULT FALSE,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
	`)
	return err
}

func (r *DB) List() ([]models.Task, error) {
	rows, err := r.db.Query("SELECT id, title, completed, created_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			log.Printf("Error scanning task row: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *DB) Create(title string, completed bool) (int64, error) {
	result, err := r.db.Exec("INSERT INTO tasks (title, completed) VALUES (?, ?)", title, completed)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *DB) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}
