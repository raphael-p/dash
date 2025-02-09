package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Importance int

const (
	Low Importance = iota
	Medium
	High
)

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
			updated_at DATETIME NOT NULL
		);`,
	down: `DROP TABLE IF EXISTS tasks`,
}

var createMetaTable = DBInitOperation{
	name: "create meta table",
	up: `
		CREATE TABLE IF NOT EXISTS meta (
			task_id INTEGER PRIMARY KEY,
			importance INTEGER NOT NULL CHECK(importance BETWEEN 0 AND 2),
			due_date DATE,
			completed_at DATETIME,
			FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
		);`,
	down: `DROP TABLE IF EXISTS meta`,
}

var initOperations = []DBInitOperation{
	createTasks,
	createMetaTable,
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
