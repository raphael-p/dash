package actions

import (
	"fmt"
	"log"
	"time"

	"github.com/raphael-p/datashard/database"
)

func AddTask(name, description string) {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	insertTask := `INSERT INTO tasks (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)`
	createdAt := time.Now()
	updatedAt := createdAt

	res, err := database.ExecWithLazyInit(tx, insertTask, name, description, createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to insert into tasks: %v", err)
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to retrieve last insert ID: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}

	fmt.Printf("task '%s' added successfully with ID %d\n", name, taskID)
}
