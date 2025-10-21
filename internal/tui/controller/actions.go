package controller

import (
	"fmt"
	"strconv"

	"github.com/raphael-p/datashard/internal/database"
)

func (c *Controller) refreshTasks() {
	err := c.displayPanel.GetCurrentPage()
	if err != nil {
		c.infoPanel.Error(err)
	}
}

func (c *Controller) openTask(idString string) {
	id, err := extractIDFromString(idString)
	if err != nil {
		c.infoPanel.Warn(fmt.Sprint("your input is invalid: ", err))
		return
	}

	task, err := c.displayPanel.ShowTaskById(int64(id))
	if err != nil {
		c.infoPanel.Error(err)
		return
	}
	c.editTaskForm(*task)
}

func (c *Controller) addTask(name, description string) {
	if name == "" || description == "" {
		c.infoPanel.Warn("please provide a name and description of the task.")
		return
	}

	task, err := database.CreateTask(name, description)
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to create task: %s", err))
		return
	}

	c.Home()
	c.infoPanel.Info(fmt.Sprintf("new task [%d] created.", task.Id))
}

func (c *Controller) removeTask(idString string) {
	id, err := extractIDFromString(idString)
	if err != nil {
		c.infoPanel.Warn(fmt.Sprint("your input is invalid: ", err))
		return
	}

	t := database.Task{Id: int64(id)}
	deleted, err := t.Delete()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to delete task: %s", err))
		return
	}

	if deleted {
		c.infoPanel.Info(fmt.Sprintf("task [%d] deleted.", id))
		c.refreshTasks()
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", id))
	}
}

func (c *Controller) editTask(task database.Task, name, description string) {
	if name == "" && description == "" {
		c.infoPanel.Warn("please update the name or description of the task.")
		return
	}

	if name != "" {
		task.Name = name
	}
	if description != "" {
		task.Description = description
	}

	updated, err := task.Update()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to edit task: %s", err))
		return
	}

	c.Home()

	if updated {
		c.infoPanel.Info(fmt.Sprintf("task [%d] updated.", task.Id))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", task.Id))
	}
}

func (c *Controller) bumpTask(idString string) {
	id, err := extractIDFromString(idString)
	if err != nil {
		c.infoPanel.Warn(fmt.Sprint("your input is invalid: ", err))
		return
	}

	bumped, err := database.BumpTask(int64(id))
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to bump task priority: %s", err))
		return
	}

	if bumped {
		c.infoPanel.Info(fmt.Sprintf("bumped priority of task [%d]", id))
		c.refreshTasks()
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", id))
	}
}

func extractIDFromString(idString string) (int, error) {
	if idString == "" {
		return 0, fmt.Errorf("task ID is empty")
	}

	id, err := strconv.Atoi(idString)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("task ID is not an integer")
	}

	if id < 1 {
		return 0, fmt.Errorf("task ID must be greater than one (1)")
	}
	return id, nil
}
