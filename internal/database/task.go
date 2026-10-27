package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/raphael-p/dash/pkg/logger"
)

type TaskMode uint8

const (
	ActiveTasks TaskMode = iota
	CompletedTasks
	PausedTasks
)

func (tm TaskMode) GetNext() TaskMode {
	return TaskMode((tm + 1) % 3)
}

func (tm TaskMode) ToString() string {
	switch tm {
	case CompletedTasks:
		return "Completed"
	case PausedTasks:
		return "Paused"
	default:
		return "Active"
	}
}

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
	PausedAt         sql.NullTime  `json:"paused_at"`
	PriorityBumpedAt sql.NullTime  `json:"priority_bumped_at"`
	TimeSpentSeconds sql.NullInt16 `json:"time_spent_seconds"`
}

func (Task) SQLColumns() string {
	var names []string
	for field := range reflect.TypeFor[Task]().Fields() {
		names = append(names, field.Tag.Get("json"))
	}
	return strings.Join(names, ", ")
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
		return Task{}, fmt.Errorf("%w: %w", errLastInsertFetchFailed, err)
	}

	return Task{taskID, name, description, createdAt, updatedAt, sql.NullTime{}, sql.NullTime{}, sql.NullTime{}, sql.NullInt16{}}, nil
}

func GetTask(id int64) (Task, error) {
	if id <= 0 {
		return Task{}, fmt.Errorf("%w, got %d", ErrIDNotPositive, id)
	}

	query := fmt.Sprintf(`SELECT %s FROM tasks WHERE id = ?;`, Task{}.SQLColumns())
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

func getTaskQueryParts(mode TaskMode, cursor TaskCursor) (where string, args []any, orderBy string) {
	var baseWhere, cursorWhere string
	switch mode {
	case CompletedTasks:
		baseWhere, cursorWhere, orderBy = "completed_at IS NOT NULL", "completed_at < ?", "completed_at DESC"
	case PausedTasks:
		baseWhere = "completed_at IS NULL AND paused_at IS NOT NULL"
		cursorWhere, orderBy = "paused_at > ?", "paused_at ASC"
	case ActiveTasks:
		baseWhere = "completed_at IS NULL AND paused_at IS NULL"
		cursorWhere = "(priority_bumped_at IS NULL OR priority_bumped_at < ?)"
		orderBy = "priority_bumped_at DESC"
	}

	if cursor.Timestamp.Valid {
		return baseWhere + " AND " + cursorWhere, []any{cursor.Timestamp}, orderBy
	} else {
		return baseWhere + " AND id > ?", []any{cursor.ID}, orderBy
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

func PauseTask(taskID int64) (bool, error) {
	pauseTask := `
	UPDATE tasks
	SET paused_at = ?
	WHERE id = ?;
	`

	logger.Debugf("pausing task [%d]", taskID)
	res, err := DB.Exec(pauseTask, time.Now(), taskID)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	return count > 0, err
}

func UnpauseTask(taskID int64) (bool, error) {
	unpauseTask := `
	UPDATE tasks
	SET paused_at = NULL
	WHERE id = ?;
	`

	logger.Debugf("unpausing task [%d]", taskID)
	res, err := DB.Exec(unpauseTask, taskID)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	return count > 0, err
}

func ReactivateTask(taskID int64) (bool, error) {
	reactivateTask := `
	UPDATE tasks
	SET completed_at = NULL
	WHERE id = ?;
	`

	logger.Debugf("reactivating task [%d]", taskID)
	res, err := DB.Exec(reactivateTask, taskID)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	return count > 0, err
}

func (t *Task) GetFormattedTimes() (createdAt, updatedAt, completedAt, pausedAt string) {
	createdAt = t.CreatedAt.Format(time.DateTime)
	updatedAt = t.UpdatedAt.Format(time.DateTime)
	if t.CompletedAt.Valid {
		completedAt = t.CompletedAt.Time.Format(time.DateTime)
	}
	if t.PausedAt.Valid {
		pausedAt = t.PausedAt.Time.Format(time.DateTime)
	}
	return createdAt, updatedAt, completedAt, pausedAt
}
