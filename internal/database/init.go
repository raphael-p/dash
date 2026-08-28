package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/raphael-p/datashard/pkg/logger"
)

const schemaVersion int = 1

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
			priority_bumped_at DATETIME,
			time_spent_seconds INTEGER
		);`,
	down: `DROP TABLE IF EXISTS tasks;`,
}

var createTaskValidationTriggers = DBInitOperation{
	name: "create task validation triggers",
	up: `
		DROP TRIGGER IF EXISTS tasks_completed_at_not_future_insert;
		DROP TRIGGER IF EXISTS tasks_completed_at_not_future_update;

		CREATE TRIGGER IF NOT EXISTS tasks_completed_at_not_future_insert
		BEFORE INSERT ON tasks
		WHEN NEW.completed_at IS NOT NULL
		 AND julianday(NEW.completed_at) > julianday('now', '+5 second')
		BEGIN
			SELECT RAISE(ABORT, 'completed_at cannot be in the future');
		END;

		CREATE TRIGGER IF NOT EXISTS tasks_completed_at_not_future_update
		BEFORE UPDATE OF completed_at ON tasks
		WHEN NEW.completed_at IS NOT NULL
		 AND julianday(NEW.completed_at) > julianday('now', '+5 second')
		BEGIN
			SELECT RAISE(ABORT, 'completed_at cannot be in the future');
		END;`,
	down: `
		DROP TRIGGER IF EXISTS tasks_completed_at_not_future_insert;
		DROP TRIGGER IF EXISTS tasks_completed_at_not_future_update;`,
}

var setSchemaVersion = DBInitOperation{
	name: "set schema version",
	up:   fmt.Sprintf("PRAGMA user_version = %d;", schemaVersion),
	down: "PRAGMA user_version = 0;",
}

var initOperations = []DBInitOperation{
	createTasks,
	createTaskValidationTriggers,
	setSchemaVersion,
}

func Initialise() {
	for _, op := range initOperations {
		_, err := DB.Exec(op.up)
		if err != nil {
			logger.Fatalf("database initialisation error on '%s': %v", op.name, err)
		}
	}
	logger.Infof("database initialised successfully with schema version %d", schemaVersion)
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

func CheckVersion() {
	var dbVersion int
	err := DB.QueryRow("PRAGMA user_version;").Scan(&dbVersion)
	if err != nil {
		logger.Fatalf("failed to read database version: %v", err)
	}

	if dbVersion > 0 {
		if dbVersion != schemaVersion {
			logger.Fatalf("database version is out-of-sync: need %d, got %d", schemaVersion, dbVersion)
		}
		logger.Debugf("database has schema version %d", dbVersion)
	}
}
