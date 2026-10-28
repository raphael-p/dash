package database

import (
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/raphael-p/dash/pkg/sqlutil"
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

func (t Task) JSON() map[string]any {
	return map[string]any{
		"id": t.ID, "name": t.Name, "description": t.Description,
		"created_at": t.CreatedAt, "updated_at": t.UpdatedAt,
		"completed_at":       sqlutil.OrNil(t.CompletedAt),
		"paused_at":          sqlutil.OrNil(t.PausedAt),
		"priority_bumped_at": sqlutil.OrNil(t.PriorityBumpedAt),
		"time_spent_seconds": sqlutil.OrNil(t.TimeSpentSeconds),
	}
}

type scannable interface {
	Scan(dest ...any) error
}

func (Task) SQLColumns() string {
	var names []string
	for field := range reflect.TypeFor[Task]().Fields() {
		names = append(names, field.Tag.Get("json"))
	}
	return strings.Join(names, ", ")
}

func taskFromRow(row scannable) (Task, error) {
	var task Task
	var id int64
	var name, description string
	var createdAt, updatedAt time.Time
	var completedAt, pausedAt, priorityBumpedAt sql.NullTime
	var timeSpentSeconds sql.NullInt16

	err := row.Scan(&id, &name, &description, &createdAt, &updatedAt, &completedAt, &pausedAt, &priorityBumpedAt, &timeSpentSeconds)
	if err != nil {
		return task, err
	}
	task = Task{id, name, description, createdAt, updatedAt, completedAt, pausedAt, priorityBumpedAt, timeSpentSeconds}
	return task, nil
}
