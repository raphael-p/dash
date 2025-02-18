package actions

import (
	"fmt"
	"time"

	"github.com/raphael-p/datashard/database"
	"github.com/raphael-p/datashard/logger"
)

func AddTask(name, description string) {
	tx, err := database.DB.Begin()
	if err != nil {
		logger.Fatalf("failed to begin transaction: %v", err)
	}

	insertTask := `INSERT INTO tasks (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)`
	createdAt := time.Now()
	updatedAt := createdAt

	res, err := database.ExecWithLazyInit(tx, insertTask, name, description, createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		logger.Fatalf("failed to insert into tasks: %v", err)
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		logger.Fatalf("failed to retrieve last insert ID: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		logger.Fatalf("failed to commit transaction: %v", err)
	}

	logger.Infof("task '%s' added successfully with ID %d\n", name, taskID)
	fmt.Println(taskID)
}
