package actions

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/raphael-p/datashard/database"
)

// ListNotes retrieves and displays all notes from the database
func ListNotes() {
	query := `
    SELECT n.id, n.title, n.content, n.created_at, n.updated_at, m.urgency, m.due_date
    FROM notes n
    JOIN meta m ON n.id = m.note_id
    ORDER BY m.urgency DESC, m.due_date ASC;
    `

	rows, err := database.QueryWithLazyInit(nil, query)
	if err != nil {
		log.Fatalf("failed to query notes: %v", err)
	}
	defer rows.Close()

	fmt.Println("===== Your Notes =====")
	for rows.Next() {
		var id, urgency int
		var title, content string
		var createdAt, updatedAt time.Time
		var dueDate sql.NullTime

		err := rows.Scan(&id, &title, &content, &createdAt, &updatedAt, &urgency, &dueDate)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}

		fmt.Printf("\nID: %d\nTitle: %s\nContent: %s\nCreated At: %s\nLast Updated: %s\nUrgency: %d\n",
			id, title, content, createdAt.Format("2006-01-02"), updatedAt.Format("2006-01-02"), urgency)
		if dueDate.Valid {
			fmt.Printf("due date: %s\n", dueDate.Time.Format("2006-01-02"))
		} else {
			fmt.Println("due date: None")
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("row error: %v", err)
	}
}
