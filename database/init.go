package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
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
			completed_at DATETIME
		);`,
	down: `DROP TABLE IF EXISTS tasks`,
}

var initOperations = []DBInitOperation{
	createTasks,
}

func Initialize() {
	for _, op := range initOperations {
		_, err := DB.Exec(op.up)
		if err != nil {
			log.Fatalf("initialisation error on '%s': %v", op.name, err)
		}
	}
	log.Println("database initialized successfully")
}

func Wipe() {
	for _, op := range initOperations {
		_, err := DB.Exec(op.down)
		if err != nil {
			log.Fatalf("rollback error on '%s': %v", op.name, err)
		}
	}
	log.Println("database wiped")
	Initialize()
}
