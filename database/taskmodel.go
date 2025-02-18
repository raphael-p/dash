package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt sql.NullTime `json:"completed_at"`
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanRow(row scannable) (Task, error) {
	var task Task
	var id int64
	var name, description string
	var createdAt, updatedAt time.Time
	var completedAt sql.NullTime

	err := row.Scan(&id, &name, &description, &createdAt, &updatedAt, &completedAt)
	if err != nil {
		return task, err
	}
	task = Task{id, name, description, createdAt, updatedAt, completedAt}
	return task, nil
}

func CreateTask(name, description string) (Task, error) {
	insertTask := `
	INSERT INTO tasks (name, description, created_at, updated_at) 
	VALUES (?, ?, ?, ?)
	`

	createdAt := time.Now()
	updatedAt := createdAt
	res, err := DB.Exec(insertTask, name, description, createdAt, updatedAt)
	if err != nil {
		return Task{}, err
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		return Task{}, fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}

	return Task{taskID, name, description, createdAt, updatedAt, sql.NullTime{}}, nil
}

func GetTask(id int64) (Task, error) {
	if id <= 0 {
		return Task{}, fmt.Errorf("task id must be > 0, got %d", id)
	}

	query := `
    SELECT id, name, description, created_at, updated_at, completed_at
    FROM tasks
	WHERE id = ?;
    `
	row := DB.QueryRow(query, id)
	return scanRow(row)
}

func GetTasks(searchQuery string) ([]Task, error) {
	query := `
    SELECT id, name, description, created_at, updated_at, completed_at
    FROM tasks
	WHERE name LIKE ?
    ORDER BY id ASC;
    `

	var tasks []Task
	rows, err := DB.Query(query, "%"+searchQuery+"%")
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		task, err := scanRow(rows)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (t *Task) MarkAsDone() error {
	updateTask := `
	UPDATE tasks 
	SET updated_at = ?, completed_at = ? 
	WHERE id = ?
	`

	updatedAt := time.Now()
	completedAt := updatedAt
	_, err := DB.Exec(updateTask, updatedAt, completedAt, t.Id)
	return err
}

func (t *Task) Delete() error {
	deleteTask := `
	DELETE FROM tasks
	WHERE id = ?
	`

	_, err := DB.Exec(deleteTask, t.Id)
	return err
}
