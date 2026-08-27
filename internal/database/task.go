package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/raphael-p/datashard/pkg/logger"
)

const taskPageLimit int = 30

type TaskMode uint8

const (
	ActiveTasks TaskMode = iota
	CompletedTasks
)

type TaskCursor struct {
	ID        int64
	Timestamp sql.NullTime
}

type Task struct {
	ID               int64         `json:"id"`
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
	CompletedAt      sql.NullTime  `json:"completed_at"`
	PriorityBumpedAt sql.NullTime  `json:"priority_bumped_at"`
	TimeSpentSeconds sql.NullInt16 `json:"time_spent_seconds"`
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

	return Task{taskID, name, description, createdAt, updatedAt, sql.NullTime{}, sql.NullTime{}, sql.NullInt16{}}, nil
}

func GetTask(id int64) (Task, error) {
	if id <= 0 {
		return Task{}, fmt.Errorf("task id must be > 0, got %d", id)
	}

	query := `SELECT * FROM tasks WHERE id = ?;`
	logger.Debugf("retrieving task [%d]", id)
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
    SELECT * FROM tasks
	WHERE name LIKE ?
    ORDER BY id ASC;
    `

	logger.Trace("retrieving all tasks")
	tasks, err := getTasksInternal(query, "%"+searchQuery+"%")
	return tasks, err
}

func GetTopTask() (Task, error) {
	query := `
    SELECT * FROM tasks
	WHERE completed_at IS NULL
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

func getTasksFilters(mode TaskMode, cursor TaskCursor) (string, any) {
	if mode == CompletedTasks && !cursor.Timestamp.Valid {
		return "completed_at IS NOT NULL AND id > ?", cursor.ID
	} else if mode == CompletedTasks {
		return "completed_at IS NOT NULL AND completed_at < ?", cursor.Timestamp
	} else if cursor.Timestamp.Valid {
		return "completed_at IS NULL AND (priority_bumped_at IS NULL OR priority_bumped_at < ?)", cursor.Timestamp
	}
	return "completed_at IS NULL AND priority_bumped_at IS NULL AND id > ?", cursor.ID
}

func GetTasksPaginated(mode TaskMode, cursor TaskCursor) ([]Task, bool, error) {
	whereClause, whereArg := getTasksFilters(mode, cursor)

	var sortColumn string
	if mode == CompletedTasks {
		sortColumn = "completed_at"
	} else {
		sortColumn = "priority_bumped_at"
	}

	query := fmt.Sprintf(`
    SELECT * FROM tasks
	WHERE %s
    ORDER BY %s DESC, id ASC
	LIMIT ?;
    `, whereClause, sortColumn)

	logger.Trace(fmt.Sprintf(
		"retrieving up to %d tasks after cursor (id: %d, timestamp: %s)",
		taskPageLimit, cursor.ID, cursor.Timestamp.Time,
	))
	tasks, err := getTasksInternal(query, whereArg, taskPageLimit+1)
	if err != nil {
		return tasks, false, err
	}

	if len(tasks) > taskPageLimit {
		return tasks[:len(tasks)-1], true, nil
	}
	return tasks, false, nil
}

func (t *Task) Update() (bool, error) {
	updateTask := `
    UPDATE tasks
    SET name = :name, 
		description = :description, 
		updated_at = :updated_at,
        completed_at = :completed_at, 
		priority_bumped_at = :priority_bumped_at,
        time_spent_seconds = :time_spent_seconds
    WHERE id = :id;
    `

	logger.Debugf("updating task [%d]", t.ID)
	res, err := DB.Exec(updateTask,
		sql.Named("id", t.ID),
		sql.Named("name", t.Name),
		sql.Named("description", t.Description),
		sql.Named("updated_at", time.Now()),
		sql.Named("completed_at", t.CompletedAt),
		sql.Named("priority_bumped_at", t.PriorityBumpedAt),
		sql.Named("time_spent_seconds", t.TimeSpentSeconds),
	)
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

	logger.Debugf("deleting task [%d]", t.ID)
	res, err := DB.Exec(deleteTask, t.ID)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	return count > 0, err
}

func BumpTask(taskID int64) (bool, error) {
	bumpTask := `
	UPDATE tasks
	SET priority_bumped_at = ?
	WHERE id = ?;
	`

	logger.Debugf("bumping task [%d] priority", taskID)
	res, err := DB.Exec(bumpTask, time.Now(), taskID)
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
