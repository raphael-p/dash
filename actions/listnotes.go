package actions

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ListNotes retrieves and displays all notes from the database
func ListNotes(db *sql.DB) {
	query := `
    SELECT n.id, n.title, n.content, n.created_at, n.updated_at, m.importance_level, m.due_date
    FROM notes n
    JOIN meta m ON n.id = m.note_id
    ORDER BY m.importance_level DESC, m.due_date ASC;
    `

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to query notes: %v", err)
	}
	defer rows.Close()

	fmt.Println("===== Your Notes =====")
	for rows.Next() {
		var id, importance int
		var title, content string
		var createdAt, updatedAt time.Time
		var dueDate sql.NullTime

		err := rows.Scan(&id, &title, &content, &createdAt, &updatedAt, &importance, &dueDate)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		fmt.Printf("\nID: %d\nTitle: %s\nContent: %s\nCreated At: %s\nLast Updated: %s\nImportance: %d\n",
			id, title, content, createdAt.Format("2006-01-02"), updatedAt.Format("2006-01-02"), importance)
		if dueDate.Valid {
			fmt.Printf("Due Date: %s\n", dueDate.Time.Format("2006-01-02"))
		} else {
			fmt.Println("Due Date: None")
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Row error: %v", err)
	}
}
