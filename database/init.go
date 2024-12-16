package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase(db *sql.DB) {
	createNotesTable := `
    CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL
    );`

	createMetaTable := `
    CREATE TABLE IF NOT EXISTS meta (
        note_id INTEGER PRIMARY KEY,
        importance_level INTEGER NOT NULL CHECK(importance_level BETWEEN 1 AND 3),
        due_date DATE,
        FOREIGN KEY(note_id) REFERENCES notes(id) ON DELETE CASCADE
    );`

	_, err := db.Exec(createNotesTable)
	if err != nil {
		log.Fatalf("error creating notes table: %v", err)
	}

	_, err = db.Exec(createMetaTable)
	if err != nil {
		log.Fatalf("error creating meta table: %v", err)
	}

	log.Println("database initialized successfully.")
}
