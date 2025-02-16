package actions

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/raphael-p/datashard/database"
)

type ListMode int

const (
	ListTodo ListMode = iota
	ListDone
	ListAll
)

func ListTasks(listMode ListMode, searchQuery string) {
	query := `
    SELECT id, name, description, created_at, updated_at, completed_at
    FROM tasks
	WHERE name LIKE ?
    ORDER BY id ASC;
    `

	rows, err := database.QueryWithLazyInit(nil, query, "%"+searchQuery+"%")
	if err != nil {
		log.Fatalf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	fmt.Println("===== your tasks =====")
	for rows.Next() {
		var id int64
		var name, description string
		var createdAt, updatedAt time.Time
		var completedAt sql.NullTime

		err := rows.Scan(&id, &name, &description, &createdAt, &updatedAt, &completedAt)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}
		isDone := completedAt.Valid
		if (listMode == ListTodo && isDone) || (listMode == ListDone && !isDone) {
			continue
		}

		fmt.Printf("\nid: %d\nname: %s\n", id, name)

		if description != "" {
			fmt.Printf("description: %s\n", description)
		}

		if isDone {
			fmt.Printf("completed at: %s\n", completedAt.Time.Format(time.DateTime))
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
