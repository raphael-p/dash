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

func (c *Controller) openTask(idString string, back func(), quit func()) bool {
	id, err := extractIDFromString(idString)
	if err != nil {
		c.infoPanel.Warn(fmt.Sprint("your input is invalid: ", err))
		return false
	}

	task, err := c.displayPanel.ShowTaskById(int64(id))
	if err != nil {
		c.infoPanel.Error(err)
		return false
	}
	c.editTaskForm(&task, back, quit)

	return true
}

func (c *Controller) addTask(name, description string) bool {
	if name == "" || description == "" {
		c.infoPanel.Warn("please provide a name and description of the task.")
		return false
	}

	task, err := database.CreateTask(name, description)
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to create task: %s", err))
		return false
	}

	c.infoPanel.Info(fmt.Sprintf("new task [%d] created.", task.Id))
	return true
}

func (c *Controller) removeTask(idString string) bool {
	id, err := extractIDFromString(idString)
	if err != nil {
		c.infoPanel.Warn(fmt.Sprint("your input is invalid: ", err))
		return false
	}

	t := database.Task{Id: int64(id)}
	deleted, err := t.Delete()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to delete task: %s", err))
		return false
	}

	if deleted {
		c.infoPanel.Info(fmt.Sprintf("task [%d] deleted.", id))
		c.refreshTasks()
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", id))
	}

	return true
}

func (c *Controller) editTask(task *database.Task, name, description string) bool {
	if name == "" && description == "" {
		c.infoPanel.Warn("please update the name or description of the task.")
		return false
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
		return false
	}

	if updated {
		c.infoPanel.Info(fmt.Sprintf("task [%d] updated.", task.Id))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", task.Id))
	}
	return true
}

func (c *Controller) bumpTask(idString string) bool {
	id, err := extractIDFromString(idString)
	if err != nil {
		c.infoPanel.Warn(fmt.Sprint("your input is invalid: ", err))
		return false
	}

	bumped, err := database.BumpTask(int64(id))
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to bump task priority: %s", err))
		return false
	}

	if bumped {
		c.infoPanel.Info(fmt.Sprintf("bumped priority of task [%d]", id))
		c.refreshTasks()
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", id))
	}
	return true
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
