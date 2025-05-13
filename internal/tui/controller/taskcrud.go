package controller

import (
	"fmt"
	"strconv"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/rivo/tview"
)

func (ac *Controller) viewTask(taskIDInput *tview.InputField) {
	id, err := extractIDFromInput(taskIDInput)
	if err != nil {
		ac.infoPanel.Warn(fmt.Sprint("Please enter a valid task ID: ", err))
		return
	}

	err = ac.displayPanel.ShowTask(int64(id))
	if err != nil {
		ac.infoPanel.Error(err)
	}
}

func (ac *Controller) addTask(taskNameInput, taskDescriptionInput *tview.InputField) {
	name := taskNameInput.GetText()
	description := taskDescriptionInput.GetText()

	if name == "" || description == "" {
		ac.infoPanel.Warn("Please provide a name and description of the task.")
		return
	}

	task, err := database.CreateTask(name, description)
	if err != nil {
		ac.infoPanel.Error(fmt.Errorf("failed to create task: %s", err))
		return
	}
	ac.infoPanel.Info(fmt.Sprintf("New task [%d] created.", task.Id))
	ac.NavigateToHome()
}

func (ac *Controller) deleteTask(taskIDInput *tview.InputField) {
	id, err := extractIDFromInput(taskIDInput)
	if err != nil {
		ac.infoPanel.Warn(fmt.Sprint("Your input is invalid: ", err))
		return
	}

	t := database.Task{Id: int64(id)}
	err = t.Delete()
	if err != nil {
		ac.infoPanel.Error(fmt.Errorf("failed to delete task: %s", err))
		return
	}
	ac.infoPanel.Info(fmt.Sprintf("Task [%d] deleted.", id))
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
