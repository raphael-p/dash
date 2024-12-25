package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Urgency int

const (
	NotUrgent Urgency = iota
	SomewhatUrgent
	Urgent
	MostUrgent
)

func Initialize(db *sql.DB) {
	var urgency Urgency = 2
	if urgency == Urgent {
		fmt.Println("a")
	}

	var wg sync.WaitGroup

	createNotesTable := `
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`
	execAsync(db, createNotesTable, "error creating notes table", &wg)

	createMetaTable := `
		CREATE TABLE IF NOT EXISTS meta (
			note_id INTEGER PRIMARY KEY,
			urgency INTEGER NOT NULL CHECK(urgency BETWEEN 0 AND 3),
			due_date DATE,
			FOREIGN KEY(note_id) REFERENCES notes(id) ON DELETE CASCADE
		);`
	execAsync(db, createMetaTable, "error creating meta table", &wg)

	wg.Wait()
	log.Println("database initialized successfully.")
}
