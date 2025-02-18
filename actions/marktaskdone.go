package actions

import (
	"github.com/raphael-p/datashard/database"
	"github.com/raphael-p/datashard/logger"
)

func MarkTaskAsDone(id int64) {
	logger.Debugf("MarkTaskAsDone invoked with id: %d", id)
	task, err := database.GetTask(id)
	if err != nil {
		database.LazyInit(err)
		logger.Debugf("noop: task (id: %d) not found, nothing to mark as done\n", id)
		return
	}

	if task.CompletedAt.Valid {
		logger.Debugf("noop: task (id: %d) already marked as done\n", id)
		return
	}

	err = task.MarkAsDone()
	if err != nil {
		logger.Fatalf("failed to mark task (id: %d) as done: %v", id, err)
	}

	logger.Infof("task (id: %d) marked as done\n", id)
}
