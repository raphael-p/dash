package actions

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/raphael-p/datashard/database"
)

type listMode int

const (
	todo listMode = iota
	done
	all
)

func listTasks(listMode listMode) {
	query := `
    SELECT n.id, n.name, n.description, n.created_at, n.updated_at, m.importance, m.due_date, m.completed_at
    FROM tasks n
    JOIN meta m ON n.id = m.task_id
    ORDER BY m.importance DESC, m.due_date ASC;
    `

	rows, err := database.QueryWithLazyInit(nil, query)
	if err != nil {
		log.Fatalf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	fmt.Println("===== your tasks =====")
	for rows.Next() {
		var id, importance int
		var name, description string
		var createdAt, updatedAt time.Time
		var dueDate, completedAt sql.NullTime

		isDone := completedAt.Valid
		if (listMode == todo && isDone) || (listMode == done && !isDone) {
			continue
		}

		err := rows.Scan(&id, &name, &description, &createdAt, &updatedAt, &importance, &dueDate, &completedAt)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}

		fmt.Printf("\nid: %d\nname: %s\nimportance: %d\n", id, name, importance)

		if description != "" {
			fmt.Printf("description: %s\n", description)
		}

		if dueDate.Valid {
			fmt.Printf("due date: %s\n", dueDate.Time.Format(time.DateOnly))
		} else {
			fmt.Println("due date: None")
		}

		if isDone {
			fmt.Printf("completed At: %s\n", completedAt.Time.Format(time.DateTime))
		}

		fmt.Printf("created at: %s\nlast updated: %s\n",
			createdAt.Format(time.DateTime),
			updatedAt.Format(time.DateTime),
		)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("row error: %v", err)
	}
}

func ListTodo() {
	listTasks(todo)
}
func ListDone() {
	listTasks(done)
}
func ListAll() {
	listTasks(all)
}
