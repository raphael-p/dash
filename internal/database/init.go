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

var setSchemaVersion = DBInitOperation{
	name: "set schema version",
	up:   fmt.Sprintf("PRAGMA user_version = %d;", schemaVersion),
	down: "PRAGMA user_version = 0;",
}

var initOperations = []DBInitOperation{
	createTasks,
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
