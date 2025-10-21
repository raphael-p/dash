package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/raphael-p/datashard/pkg/logger"
)

const taskPageLimit int = 30

type Task struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt sql.NullTime `json:"completed_at"`
}

func CreateTask(name, description string) (Task, error) {
	insertTask := `
	INSERT INTO tasks (name, description, created_at, updated_at) 
	VALUES (?, ?, ?, ?)
	`

	createdAt := time.Now()
	updatedAt := createdAt
	logger.Debugf("creating task with name '%s'", name)
	res, err := DB.Exec(insertTask, name, description, createdAt, updatedAt)
	if lazyInit(err) {
		logger.Trace("retrying task creation")
		res, err = DB.Exec(insertTask, name, description, createdAt, updatedAt)
	}
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
	logger.Debugf("retrieving task (id: %d)", id)
	row := DB.QueryRow(query, id)
	task, err := scanRow(row)
	if lazyInit(err) {
		logger.Trace("retrying task retrieval")
		task, err = scanRow(row)
	}
	return task, err
}

func getTasksInternal(sqlQuery string, queryArgs ...any) ([]Task, error) {
	var tasks []Task
	rows, err := DB.Query(sqlQuery, queryArgs...)
	if lazyInit(err) {
		logger.Trace("retrying task retrieval")
		rows, err = DB.Query(sqlQuery, queryArgs...)
	}
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

func SearchTasks(searchQuery string) ([]Task, error) {
	query := `
    SELECT id, name, description, created_at, updated_at, completed_at
    FROM tasks
	WHERE name LIKE ?
    ORDER BY id ASC;
    `

	logger.Trace("retrieving all tasks")
	tasks, err := getTasksInternal(query, "%"+searchQuery+"%")
	return tasks, err
}

func GetTopTask() (Task, error) {
	query := `
    SELECT id, name, description, created_at, updated_at, completed_at
    FROM tasks
    ORDER BY priority_bumped_at DESC, id ASC
	LIMIT 1;
    `

	tasks, err := getTasksInternal(query)
	if err != nil {
		return Task{}, err
	}
	if len(tasks) == 0 {
		return Task{}, fmt.Errorf("no task found")
	}

	return tasks[0], err
}

func GetTasksPaginated(fromID int64) ([]Task, bool, error) {
	query := `
    SELECT id, name, description, created_at, updated_at, completed_at
    FROM tasks
	WHERE id > ?
    ORDER BY priority_bumped_at DESC, id ASC
	LIMIT ?;
    `

	logger.Trace(fmt.Sprintf(
		"retrieving tasks up to %d after id %d (paginated)",
		taskPageLimit,
		fromID,
	))
	tasks, err := getTasksInternal(query, fromID, taskPageLimit)
	if err != nil {
		return tasks, false, err
	}

	var hasNext bool
	if len(tasks) >= taskPageLimit {
		logger.Trace("checking if there are more tasks")

		hasNextQuery := `
		SELECT 1
		FROM tasks
		WHERE id > ?
		LIMIT 1;
		`

		lastID := tasks[len(tasks)-1].Id
		err := DB.QueryRow(hasNextQuery, lastID).Scan(&hasNext)
		if err == sql.ErrNoRows {
			hasNext = false
		} else if err != nil {
			logger.Warning(fmt.Sprint("error occured while checking if there are more tasks: ", err))
		}
	}

	return tasks, hasNext, nil
}

func BumpTask(taskID int64) (bool, error) {
	bumpTask := `
	UPDATE tasks
	SET priority_bumped_at = ?
	WHERE id = ?;
	`

	logger.Debugf("bumping task priority (id: %d)", taskID)
	res, err := DB.Exec(bumpTask, time.Now(), taskID)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	return count > 0, err
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

func (t *Task) Update() (bool, error) {
	updateTask := `
	UPDATE tasks
	SET name = ?, description = ?, updated_at = ?
	WHERE id = ?;
	`

	logger.Debugf("updating task (id: %d)", t.Id)
	res, err := DB.Exec(updateTask, t.Name, t.Description, time.Now(), t.Id)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	return count > 0, err
}

func (t *Task) Delete() (bool, error) {
	deleteTask := `
	DELETE FROM tasks
	WHERE id = ?
	`

	logger.Debugf("deleting task (id: %d)", t.Id)
	res, err := DB.Exec(deleteTask, t.Id)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	return count > 0, err
}

func (t *Task) GetFormattedTimes() (createdAt string, updatedAt string, completedAt string) {
	createdAt = t.CreatedAt.Format(time.DateTime)
	updatedAt = t.UpdatedAt.Format(time.DateTime)
	if t.CompletedAt.Valid {
		completedAt = t.CompletedAt.Time.Format(time.DateTime)
	}
	return createdAt, updatedAt, completedAt
}
