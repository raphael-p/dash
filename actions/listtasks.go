package actions

import (
	"fmt"
	"time"

	"github.com/raphael-p/datashard/database"
	"github.com/raphael-p/datashard/logger"
)

type ListMode int

const (
	ListTodo ListMode = iota
	ListDone
	ListAll
)

func ListTasks(listMode ListMode, searchQuery string) {
	tasks, err := database.GetTasks(searchQuery)
	if err != nil {
		database.LazyInit(err)
		logger.Fatalf("failed to query tasks: %v", err)
	}

	logger.Trace("printing task list")
	for _, task := range tasks {
		isDone := task.CompletedAt.Valid
		if (listMode == ListTodo && isDone) || (listMode == ListDone && !isDone) {
			continue
		}

		fmt.Printf("\nid: %d\nname: %s\n", task.Id, task.Name)

		if task.Description != "" {
			fmt.Printf("description: %s\n", task.Description)
		}

		if isDone {
			fmt.Printf("completed at: %s\n", task.CompletedAt.Time.Format(time.DateTime))
		}

		fmt.Printf("created at: %s\nlast updated: %s\n",
			task.CreatedAt.Format(time.DateTime),
			task.UpdatedAt.Format(time.DateTime),
		)
	}
}
