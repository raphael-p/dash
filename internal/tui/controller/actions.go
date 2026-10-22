package controller

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/raphael-p/dash/internal/database"
)

func (c *Controller) refreshTasks() {
	err := c.displayPanel.GetCurrentPage()
	if err != nil {
		c.infoPanel.Error(err)
	}
}

func (c *Controller) openTask(id int64, back func()) bool {
	task, err := c.displayPanel.ShowTaskByID(int64(id))
	if err != nil {
		c.infoPanel.Error(err)
		return false
	}
	c.editTaskFlow(&task, back)
	return true
}

type taskInsertionOrder int

const (
	InsertFirst taskInsertionOrder = iota + 1
	InsertNext
	InsertLast
)

func (c *Controller) addTask(name, description string, insertionOrder taskInsertionOrder) bool {
	if name == "" || description == "" {
		c.infoPanel.Warn("please provide a name and description of the task.")
		return false
	}

	task, err := database.CreateTask(name, description)
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedCreate, err))
		return false
	}

	switch insertionOrder {
	case InsertFirst:
		database.BumpTask(task.ID)
	case InsertNext:
		topTask, err := database.GetTopTask()
		if err != nil {
			c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedCreate, err))
			return false
		}
		database.BumpTask(task.ID)
		database.BumpTask(topTask.ID)
	case InsertLast: // no action required
	default:
		c.infoPanel.Error(errInvalidInsertion)
		return false
	}

	c.infoPanel.Info(fmt.Sprintf("new task [%d] created.", task.ID))
	return true
}

func (c *Controller) removeTask(id int64) bool {
	t := database.Task{ID: id}
	deleted, err := t.Delete()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedDelete, err))
		return false
	}

	if deleted {
		c.infoPanel.Info(fmt.Sprintf("task [%d] deleted.", id))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", id))
	}
	return deleted
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
		c.infoPanel.Error(fmt.Errorf("%w [%d]: %w", errFailedEdit, task.ID, err))
		return false
	}

	if updated {
		c.infoPanel.Info(fmt.Sprintf("task [%d] updated.", task.ID))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", task.ID))
	}
	return updated
}

func (c *Controller) bumpTask(id int64) bool {
	bumped, err := database.BumpTask(int64(id))
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("%w [%d]: %w", errFailedBump, id, err))
		return false
	}

	if bumped {
		c.infoPanel.Info(fmt.Sprintf("bumped priority of task [%d]", id))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", id))
	}
	return bumped
}

func (c *Controller) copyTask(id int64) {
	task, err := database.GetTask(id)
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("%w [%d]: %w", errFailedGet, id, err))
		return
	}

	err = clipboard.WriteAll(fmt.Sprintf("%s\n%s", task.Name, task.Description))
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("%w [%d]: %w", errFailedCopy, id, err))
		return
	}

	c.infoPanel.Info(fmt.Sprintf("copied task [%d] to clipboard", id))
}

func (c *Controller) reactivateTask(id int64) bool {
	reactivated, err := database.ReactivateTask(id)
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("%w [%d]: %w", errFailedReactivate, id, err))
		return false
	}

	if reactivated {
		c.infoPanel.Info(fmt.Sprintf("reactivated task [%d]", id))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop.", id))
	}
	return reactivated
}
