package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/raphael-p/dash/pkg/logger"
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

func executeTasksQuery(sqlConditions string, sqlConditionsArgs []any, sqlSort string, limit int) ([]Task, error) {
	var limitClause string
	if limit != 0 {
		limitClause = fmt.Sprint("LIMIT ", limit)
	}

	query := fmt.Sprintf(`
    SELECT * FROM tasks
	WHERE %s
    ORDER BY %s, id ASC %s;
    `, sqlConditions, sqlSort, limitClause)

	var tasks []Task
	rows, err := DB.Query(query, sqlConditionsArgs...)
	if lazyInit(err) {
		logger.Trace("retrying task retrieval")
		rows, err = DB.Query(query, sqlConditionsArgs...)
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

func GetCompletedTasksSince(since time.Time) ([]Task, error) {
	return executeTasksQuery("completed_at IS NOT NULL AND completed_at >= ?", []any{since}, "completed_at ASC", 0)
}

func getTaskQueryParts(mode TaskMode, cursor TaskCursor) (string, []any, string) {
	if mode == CompletedTasks && !cursor.Timestamp.Valid {
		return "completed_at IS NOT NULL AND id > ?", []any{cursor.ID}, "completed_at DESC"
	} else if mode == CompletedTasks {
		return "completed_at IS NOT NULL AND completed_at < ?", []any{cursor.Timestamp}, "completed_at DESC"
	} else if cursor.Timestamp.Valid {
		return "completed_at IS NULL AND (priority_bumped_at IS NULL OR priority_bumped_at < ?)", []any{cursor.Timestamp}, "priority_bumped_at DESC"
	} else {
		return "completed_at IS NULL AND id > ?", []any{cursor.ID}, "priority_bumped_at DESC"
	}
}

func GetTopTask() (Task, error) {
	sqlConditions, sqlConditionsArgs, sortColumn := getTaskQueryParts(ActiveTasks, TaskCursor{})
	tasks, err := executeTasksQuery(sqlConditions, sqlConditionsArgs, sortColumn, 1)
	if err != nil {
		return Task{}, err
	}
	if len(tasks) == 0 {
		return Task{}, fmt.Errorf("no task found")
	}

	return tasks[0], err
}

func GetTasksPaginated(mode TaskMode, cursor TaskCursor) ([]Task, bool, error) {
	logger.Trace(fmt.Sprintf(
		"retrieving up to %d tasks after cursor (id: %d, timestamp: %s)",
		taskPageLimit, cursor.ID, cursor.Timestamp.Time,
	))

	sqlConditions, sqlConditionsArgs, sortColumn := getTaskQueryParts(mode, cursor)
	tasks, err := executeTasksQuery(sqlConditions, sqlConditionsArgs, sortColumn, taskPageLimit+1)
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
