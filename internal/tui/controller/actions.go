package controller

import (
	"fmt"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/taskform"
)

func (c *Controller) refreshTasks() {
	err := c.displayPanel.GetCurrentPage()
	if err != nil {
		c.infoPanel.Error(err)
	}
}

func (c *Controller) openTask(idString string, back func()) bool {
	id, ok := getIDOrWarn(idString, c.infoPanel.Warn)
	if !ok {
		return false
	}

	task, err := c.displayPanel.ShowTaskByID(int64(id))
	if err != nil {
		c.infoPanel.Error(err)
		return false
	}
	c.editTaskForm(&task, back)

	return true
}

func (c *Controller) addTask(name, description string, priority taskform.TaskPriority) bool {
	if name == "" || description == "" {
		c.infoPanel.Warn("please provide a name and description of the task.")
		return false
	}

	task, err := database.CreateTask(name, description)
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to create task: %s", err))
		return false
	}

	switch priority {
	case taskform.PriorityFirst:
		database.BumpTask(task.ID)
	case taskform.PriorityNext:
		topTask, err := database.GetTopTask()
		if err != nil {
			c.infoPanel.Error(fmt.Errorf("failed to create task: %s", err))
			return false
		}
		database.BumpTask(task.ID)
		database.BumpTask(topTask.ID)
	case taskform.PriorityLast: // no action required
	}

	c.infoPanel.Info(fmt.Sprintf("new task [%d] created.", task.ID))
	return true
}

func (c *Controller) removeTask(idString string) bool {
	id, ok := getIDOrWarn(idString, c.infoPanel.Warn)
	if !ok {
		return false
	}

	t := database.Task{ID: int64(id)}
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
		c.infoPanel.Error(fmt.Errorf("failed to edit task [%d]: %s", task.ID, err))
		return false
	}

	if updated {
		c.infoPanel.Info(fmt.Sprintf("task [%d] updated.", task.ID))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", task.ID))
	}
	return true
}

func (c *Controller) bumpTask(idString string) bool {
	id, ok := getIDOrWarn(idString, c.infoPanel.Warn)
	if !ok {
		return false
	}

	bumped, err := database.BumpTask(int64(id))
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to bump task [%d] priority: %s", id, err))
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

func (c *Controller) copyTask(idString string) {
	id, ok := getIDOrWarn(idString, c.infoPanel.Warn)
	if !ok {
		return
	}

	task, err := database.GetTask(int64(id))
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("could not retrieve task [%d]: %s", id, err))
		return
	}

	err = clipboard.WriteAll(fmt.Sprintf("%s\n%s", task.Name, task.Description))
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to copy task [%d] to clipboard: %s", id, err))
		return
	}

	c.infoPanel.Info(fmt.Sprintf("copied task [%d] to clipboard", id))
}

func getIDOrWarn(idString string, warn func(err string)) (int, bool) {
	id, err := extractIDFromString(idString)
	if err != nil {
		warn(fmt.Sprint("your input is invalid: ", err))
		return 0, false
	}
	return id, true
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
