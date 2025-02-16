package actions

import (
	"fmt"
	"log"

	"github.com/raphael-p/datashard/database"
)

func MarkTaskAsDone(id int64) {
	task, err := database.GetTask(id)
	if err != nil {
		fmt.Printf("noop, task (id: %d) not found, nothing to mark as done\n", id)
		return
	}

	if task.CompletedAt.Valid {
		fmt.Printf("noop, task (id: %d) already marked as done\n", id)
		return
	}

	err = task.MarkAsDone()
	if err != nil {
		log.Fatalf("failed to mark task (id: %d) as done: %v", id, err)
	}

	fmt.Printf("task (id: %d) marked as done\n", id)
}
