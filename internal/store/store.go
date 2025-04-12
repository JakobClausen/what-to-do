package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type Todo struct {
	ID          int64
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
}

type TodoStore struct {
	db *sql.DB
}

func CreateStore() (*TodoStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(homeDir, ".what-to-do", "store")

	if err = os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(appDir, "store.sqlite")
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        completed BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
	`); err != nil {
		return nil, err
	}

	return &TodoStore{db: db}, nil
}

func (s *TodoStore) Close() error {
	return s.db.Close()
}

func (s *TodoStore) CreateTodo(title, description string) (*Todo, error) {
	result, err := s.db.Exec(
		"INSERT INTO todos (title, description) VALUES (?, ?)",
		title, description,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetTodo(id)
}

func (s *TodoStore) GetTodo(id int64) (*Todo, error) {
	var todo Todo
	err := s.db.QueryRow(
		"SELECT id, title, description, completed, created_at FROM todos WHERE id = ?",
		id,
	).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodoStore) ListTodos() ([]*Todo, error) {
	rows, err := s.db.Query("SELECT id, title, description, completed, created_at FROM todos ORDER BY created_at DESC")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, rows.Err()
}

func (s *TodoStore) UpdateTodo(todo *Todo) error {
	_, err := s.db.Exec(
		"UPDATE todos SET title = ?, description = ?, completed = ? WHERE id = ?",
		todo.Title, todo.Description, todo.Completed, todo.ID,
	)
	return err
}

func (s *TodoStore) DeleteTodo(id int64) error {
	_, err := s.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
