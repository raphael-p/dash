package actions

import (
	"fmt"
	"log"
	"time"

	"github.com/raphael-p/datashard/database"
)

// AddNote inserts a new note and its metadata into the database
func AddNote(title, content string, urgency int, dueDate string) {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	// Insert into notes table
	insertNote := `INSERT INTO notes (title, content, created_at, updated_at) VALUES (?, ?, ?, ?)`
	createdAt := time.Now()
	updatedAt := createdAt

	res, err := database.ExecWithLazyInit(tx, insertNote, title, content, createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to insert into notes: %v", err)
	}

	noteID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to retrieve last insert ID: %v", err)
	}

	// Insert into meta table
	insertMeta := `INSERT INTO meta (note_id, urgency, due_date) VALUES (?, ?, ?)`
	var dueDateParsed interface{}
	if dueDate != "" {
		dueDateParsed, err = time.Parse("2006-01-02", dueDate)
		if err != nil {
			tx.Rollback()
			log.Fatalf("unvalid due date format, use YYYY-MM-DD")
		}
	} else {
		dueDateParsed = nil
	}

	_, err = tx.Exec(insertMeta, noteID, urgency, dueDateParsed)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to insert into meta: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}

	fmt.Printf("note '%s' added successfully with ID %d\n", title, noteID)
}
