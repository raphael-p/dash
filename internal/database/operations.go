package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/raphael-p/dash/pkg/logger"
)

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
		return Task{}, fmt.Errorf("%w: %w", errLastInsertFetchFailed, err)
	}

	return Task{taskID, name, description, createdAt, updatedAt, sql.NullTime{}, sql.NullTime{}, sql.NullTime{}, sql.NullInt16{}}, nil
}

func GetTask(id int64) (Task, error) {
	if id <= 0 {
		return Task{}, fmt.Errorf("%w, got %d", errIDNotPositive, id)
	}

	query := fmt.Sprintf(`SELECT %s FROM tasks WHERE id = ?;`, Task{}.SQLColumns())
	logger.Debugf("retrieving task [%d]", id)
	row := DB.QueryRow(query, id)
	task, err := taskFromRow(row)
	if lazyInit(err) {
		logger.Trace("retrying task retrieval")
		task, err = taskFromRow(row)
	}
	return task, err
}

func executeTasksQuery(sqlConditions string, sqlConditionsArgs []any, sqlSort string, limit int) ([]Task, error) {
	var limitClause string
	if limit != 0 {
		limitClause = fmt.Sprint("LIMIT ", limit)
	}

	query := fmt.Sprintf(`
    SELECT %s FROM tasks
	WHERE %s
    ORDER BY %s, id ASC %s;
    `, Task{}.SQLColumns(), sqlConditions, sqlSort, limitClause)

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
		task, err := taskFromRow(rows)
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

func getTaskQueryParts(mode TaskMode, cursor TaskCursor) (where string, args []any, orderBy string) {
	var whereBase, whereID, whereTS string
	switch mode {
	case CompletedTasks:
		orderBy = "completed_at DESC"
		whereBase, whereTS = "completed_at IS NOT NULL", "completed_at < ?"
	case PausedTasks:
		orderBy = "paused_at ASC"
		whereBase, whereTS = "completed_at IS NULL AND paused_at IS NOT NULL", "paused_at > ?"
	case ActiveTasks:
		orderBy = "priority_bumped_at DESC"
		whereBase = "completed_at IS NULL AND paused_at IS NULL"
		whereID = "priority_bumped_at IS NULL AND ID > ?"
		whereTS = "priority_bumped_at IS NULL OR priority_bumped_at < ?"
	}

	if cursor.Timestamp.Valid {
		where = fmt.Sprintf("(%s) AND (%s)", whereBase, whereTS)
		return where, []any{cursor.Timestamp}, orderBy
	} else if cursor.ID > 0 && whereID != "" {
		where = fmt.Sprintf("(%s) AND (%s)", whereBase, whereID)
		return where, []any{cursor.ID}, orderBy
	} else {
		return whereBase, []any{}, orderBy
	}
}

func GetTopTask() (Task, error) {
	sqlConditions, sqlConditionsArgs, sortColumn := getTaskQueryParts(ActiveTasks, TaskCursor{})
	tasks, err := executeTasksQuery(sqlConditions, sqlConditionsArgs, sortColumn, 1)
	if err != nil {
		return Task{}, err
	}
	if len(tasks) == 0 {
		return Task{}, ErrNotFound
	}

	return tasks[0], err
}

func GetTasksPaginated(mode TaskMode, cursor TaskCursor, taskPageLimit int) ([]Task, bool, error) {
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
		paused_at = :paused_at,
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
		sql.Named("paused_at", t.PausedAt),
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

func lazyInit(err error) bool {
	if err != nil && strings.Contains(err.Error(), "no such table") {
		logger.Warning("no database found, intialising")
		Initialise()
		return true
	}
	return false
}
