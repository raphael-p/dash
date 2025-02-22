package actions

import (
	"fmt"

	"github.com/raphael-p/datashard/database"
	"github.com/raphael-p/datashard/logger"
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
