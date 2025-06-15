package controller

import (
	"fmt"
	"strconv"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/rivo/tview"
)

func (c *Controller) submitViewTask(taskIDInput *tview.InputField) {
	id, err := extractIDFromInput(taskIDInput)
	if err != nil {
		c.infoPanel.Warn(fmt.Sprint("Your input is invalid: ", err))
		return
	}

	err = c.displayPanel.ShowTask(int64(id))
	if err != nil {
		c.infoPanel.Error(err)
		return
	}

	taskIDInput.SetText("")
}

func (c *Controller) submitAddTask(taskNameInput, taskDescriptionInput *tview.InputField) {
	name := taskNameInput.GetText()
	description := taskDescriptionInput.GetText()
	if name == "" || description == "" {
		c.infoPanel.Warn("Please provide a name and description of the task.")
		return
	}

	task, err := database.CreateTask(name, description)
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to create task: %s", err))
		return
	}

	taskNameInput.SetText("")
	taskDescriptionInput.SetText("")
	c.refreshTasks()
	c.infoPanel.Info(fmt.Sprintf("New task [%d] created.", task.Id))
}

func (c *Controller) submitDeleteTask(taskIDInput *tview.InputField) {
	id, err := extractIDFromInput(taskIDInput)
	if err != nil {
		c.infoPanel.Warn(fmt.Sprint("Your input is invalid: ", err))
		return
	}

	t := database.Task{Id: int64(id)}
	deleted, err := t.Delete()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to delete task: %s", err))
		return
	}

	taskIDInput.SetText("")
	c.refreshTasks()

	if deleted {
		c.infoPanel.Info(fmt.Sprintf("Task [%d] deleted.", id))
	} else {
		c.infoPanel.Warn(fmt.Sprintf("Task [%d] does not exist, noop.", id))
	}
}

func extractIDFromInput(input *tview.InputField) (int, error) {
	idStr := input.GetText()
	if idStr == "" {
		return 0, fmt.Errorf("task ID is empty")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("task ID is not an integer")
	}

	if id < 1 {
		return 0, fmt.Errorf("task ID must be greater than one (1)")
	}
	return id, nil
}
