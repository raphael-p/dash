package controller

import (
	"fmt"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/sqlutil"
	"github.com/raphael-p/dash/pkg/strutil"
)

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
		task.PriorityBumpedAt = sqlutil.NowNullTime()
		task.Update()
	case InsertNext:
		topTask, err := database.GetTopTask()
		if err != nil {
			c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedCreate, err))
			return false
		}
		task.PriorityBumpedAt = sqlutil.NowNullTime()
		task.Update()
		topTask.PriorityBumpedAt = sqlutil.NowNullTime()
		topTask.Update()
	case InsertLast: // no action required
	default:
		c.infoPanel.Error(errInvalidInsertion)
		return false
	}

	c.infoPanel.Info(fmt.Sprintf("new task [%d] created.", task.ID))
	return true
}

func (c *Controller) getTask(id int64) (database.Task, bool) {
	task, err := database.GetTask(int64(id))
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("%w [%d]: %w", errFailedGet, id, err))
		return task, false

	}
	return task, true
}

func (c *Controller) getTaskFromCommandLine() *database.Task {
	if id, ok := strutil.ParseIDOrWarn(c.commandLine.GetInput(), c.infoPanel.Warn); ok {
		if task, ok := c.getTask(id); ok {
			return &task
		}
	}
	return nil
}

func (c *Controller) openTask(id int64, back func()) bool {
	task, ok := c.getTask(id)
	if !ok {
		return false
	}
	task, err := c.displayPanel.ShowTask(task)
	if err != nil {
		c.infoPanel.Error(err)
		return false
	}
	c.editTaskFlow(&task, back)
	return true
}

func (c *Controller) refreshTasks() {
	err := c.displayPanel.GetCurrentPage()
	if err != nil {
		c.infoPanel.Error(err)
	}
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

func (c *Controller) updateTask(task *database.Task, operation string) bool {
	updated, err := task.Update()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("%w: %s on task [%d]: %w", errFailedUpdate, operation, task.ID, err))
		return false
	}

	if updated {
		c.infoPanel.Info(fmt.Sprintf("task [%d]: %s success", task.ID, operation))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("task [%d] does not exist, noop", task.ID))
	}
	return updated
}
