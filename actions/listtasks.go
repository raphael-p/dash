package actions

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/raphael-p/datashard/database"
)

// ListTasks retrieves and displays all tasks from the database
func ListTasks() {
	query := `
    SELECT n.id, n.name, n.description, n.created_at, n.updated_at, m.importance, m.due_date
    FROM tasks n
    JOIN meta m ON n.id = m.task_id
    ORDER BY m.importance DESC, m.due_date ASC;
    `

	rows, err := database.QueryWithLazyInit(nil, query)
	if err != nil {
		log.Fatalf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	fmt.Println("===== Your Tasks =====")
	for rows.Next() {
		var id, importance int
		var name, description string
		var createdAt, updatedAt time.Time
		var dueDate sql.NullTime

		err := rows.Scan(&id, &name, &description, &createdAt, &updatedAt, &importance, &dueDate)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}

		fmt.Printf("\nID: %d\nName: %s\nDescription: %s\nCreated At: %s\nLast Updated: %s\nImportance: %d\n",
			id, name, description, createdAt.Format(time.DateTime), updatedAt.Format(time.DateTime), importance)
		if dueDate.Valid {
			fmt.Printf("due date: %s\n", dueDate.Time.Format(time.DateOnly))
		} else {
			fmt.Println("due date: None")
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("row error: %v", err)
	}
}
