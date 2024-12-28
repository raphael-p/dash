package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Urgency int

const (
	NotUrgent Urgency = iota
	SomewhatUrgent
	Urgent
	MostUrgent
)

type DBInitOperation struct {
	name,
	up,
	down string
}

var createNotes = DBInitOperation{
	name: "create notes table",
	up: `
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
	down: `DROP TABLE IF EXISTS notes`,
}

var createMetaTable = DBInitOperation{
	name: "create meta table",
	up: `
		CREATE TABLE IF NOT EXISTS meta (
			note_id INTEGER PRIMARY KEY,
			urgency INTEGER NOT NULL CHECK(urgency BETWEEN 0 AND 3),
			due_date DATE,
			FOREIGN KEY(note_id) REFERENCES notes(id) ON DELETE CASCADE
		);`,
	down: `DROP TABLE IF EXISTS meta`,
}

var initOperations = []DBInitOperation{
	createNotes,
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
