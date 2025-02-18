package actions

import (
	"github.com/raphael-p/datashard/database"
	"github.com/raphael-p/datashard/logger"
)

func DeleteTask(id int64) {
	logger.Debugf("DeleteTask invoked with id: %d", id)
	task, err := database.GetTask(id)
	if err != nil {
		database.LazyInit(err)
		logger.Debugf("noop: task (id: %d) not found, nothing to delete\n", id)
		return
	}

	if task.CompletedAt.Valid {
		logger.Debugf("noop; task (id: %d) is already done, cannot be deleted\n", id)
		return
	}

	err = task.Delete()
	if err != nil {
		logger.Debugf("failed delete task (id: %d): %v", id, err)
	}

	logger.Infof("task (id: %d) deleted\n", id)
}
