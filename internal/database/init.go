package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/raphael-p/datashard/pkg/logger"
)

var DB *sql.DB

type DBInitOperation struct {
	name,
	up,
	down string
}

var createTasks = DBInitOperation{
	name: "create tasks table",
	up: `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			completed_at DATETIME,
			priority_bumped_at DATETIME
		);`,
	down: `DROP TABLE IF EXISTS tasks`,
}

var initOperations = []DBInitOperation{
	createTasks,
}

func Initialise() {
	for _, op := range initOperations {
		_, err := DB.Exec(op.up)
		if err != nil {
			logger.Fatalf("database initialisation error on '%s': %v", op.name, err)
		}
	}
	logger.Info("database initialised successfully")
}

func Wipe() {
	for _, op := range initOperations {
		_, err := DB.Exec(op.down)
		if err != nil {
			logger.Fatalf("rollback error on '%s': %v", op.name, err)
		}
	}
	logger.Info("database wiped")
	Initialise()
}
