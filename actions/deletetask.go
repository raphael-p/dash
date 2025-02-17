package actions

import (
	"fmt"
	"log"

	"github.com/raphael-p/datashard/database"
)

func DeleteTask(id int64) {
	task, err := database.GetTask(id)
	if err != nil {
		database.LazyInit(err)
		fmt.Printf("noop: task (id: %d) not found, nothing to delete\n", id)
		return
	}

	if task.CompletedAt.Valid {
		fmt.Printf("noop; task (id: %d) is already done, cannot be deleted\n", id)
		return
	}

	err = task.Delete()
	if err != nil {
		log.Fatalf("failed delete task (id: %d): %v", id, err)
	}

	fmt.Printf("task (id: %d) deleted\n", id)
}
