package database

import (
	"database/sql"
	"strings"
	"time"

	"github.com/raphael-p/datashard/pkg/logger"
)

func lazyInit(err error) bool {
	if err != nil && strings.Contains(err.Error(), "no such table") {
		logger.Warning("no database found, intialising")
		Initialise()
		return true
	}
	return false
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
