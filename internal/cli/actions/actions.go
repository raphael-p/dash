package actions

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/logger"
)

func AddTask(name, description string) {
	logger.Debugf("AddTask invoked with name: %s, description: %s", name, description)
	task, err := database.CreateTask(name, description)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Infof("task '%s' added successfully with ID %d\n", name, task.Id)
	fmt.Println(task.Id)
}

func DeleteTask(id int64) {
	logger.Debugf("DeleteTask invoked with id: %d", id)
	task, err := database.GetTask(id)
	if err != nil {
		logger.Debugf("noop: task (id: %d) not found, nothing to delete\n", id)
		return
	}

	if task.CompletedAt.Valid {
		logger.Debugf("noop; task (id: %d) is already done, cannot be deleted\n", id)
		return
	}

	deleted, err := task.Delete()
	if err != nil {
		logger.Debugf("failed delete task (id: %d): %v", id, err)
	}

	if deleted {
		logger.Infof("task (id: %d) deleted\n", id)
	} else {
		logger.Warningf("task (id: %d) does not exist, noop\n", id)
	}
}

type ListMode int

const (
	ListTodo ListMode = iota
	ListDone
	ListAll
)

func ListTasks(listMode ListMode, searchQuery string) {
	logger.Debugf("ListTasks invoked with listMode: %d, searchQuery: %s", listMode, searchQuery)
	tasks, err := database.SearchTasks(searchQuery)
	if err != nil {
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

func MarkTaskAsDone(id int64) {
	logger.Debugf("MarkTaskAsDone invoked with id: %d", id)
	task, err := database.GetTask(id)
	if err != nil {
		logger.Debugf("noop: task (id: %d) not found, nothing to mark as done\n", id)
		return
	}

	if task.CompletedAt.Valid {
		logger.Debugf("noop: task (id: %d) already marked as done\n", id)
		return
	}

	task.CompletedAt = sql.NullTime{Time: time.Now()}
	_, err = task.Update()
	if err != nil {
		logger.Fatalf("failed to mark task (id: %d) as done: %v", id, err)
	}

	logger.Infof("task (id: %d) marked as done\n", id)
}
