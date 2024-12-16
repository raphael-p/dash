package actions

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// AddNote inserts a new note and its metadata into the database
func AddNote(db *sql.DB, title, content string, importance int, dueDate string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	// Insert into notes table
	insertNote := `INSERT INTO notes (title, content, created_at, updated_at) VALUES (?, ?, ?, ?)`
	createdAt := time.Now()
	updatedAt := createdAt

	res, err := tx.Exec(insertNote, title, content, createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert into notes: %v", err)
	}

	noteID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to retrieve last insert ID: %v", err)
	}

	// Insert into meta table
	insertMeta := `INSERT INTO meta (note_id, importance_level, due_date) VALUES (?, ?, ?)`
	var dueDateParsed interface{}
	if dueDate != "" {
		dueDateParsed, err = time.Parse("2006-01-02", dueDate)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Invalid due date format. Use YYYY-MM-DD.")
		}
	} else {
		dueDateParsed = nil
	}

	_, err = tx.Exec(insertMeta, noteID, importance, dueDateParsed)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert into meta: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Printf("Note '%s' added successfully with ID %d.\n", title, noteID)
}
