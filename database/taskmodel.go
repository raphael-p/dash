package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	Id                   int64
	Name, Description    string
	CreatedAt, UpdatedAt time.Time
	CompletedAt          sql.NullTime
}

func GetTask(id int64) (Task, error) {
	query := `
    SELECT name, description, created_at, updated_at, completed_at
    FROM tasks
	WHERE id = ?;
    `

	var task Task
	if id <= 0 {
		return task, fmt.Errorf("task id must be > 0, got %d", id)
	}

	var name, description string
	var createdAt, updatedAt time.Time
	var completedAt sql.NullTime
	row := DB.QueryRow(query, id)
	err := row.Scan(&name, &description, &createdAt, &updatedAt, &completedAt)
	if err != nil {
		return task, err
	}

	task = Task{id, name, description, createdAt, updatedAt, completedAt}
	return task, nil
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
