package tui

import (
	"fmt"
	"strconv"

	"github.com/raphael-p/datashard/database"
	"github.com/rivo/tview"
)

func (ac *appController) viewTaskFunc(taskIDInput *tview.InputField) func() {
	return func() {
		id, err := extractIDFromInput(taskIDInput)
		if err != nil {
			ac.infoPanel.error(fmt.Sprint("Your input is invalid: ", err))
			return
		}

		err = ac.displayPanel.showTask(int64(id))
		if err != nil {
			ac.infoPanel.error(err)
		}
	}
}

func (ac *appController) addTaskFunc(taskNameInput, taskDescriptionInput *tview.InputField) func() {
	return func() {
		name := taskNameInput.GetText()
		description := taskDescriptionInput.GetText()

		if name == "" || description == "" {
			ac.infoPanel.warning("Please provide a name and description of the task.")
			return
		}

		task, err := database.CreateTask(name, description)
		if err != nil {
			ac.infoPanel.error(err)
			return
		}
		ac.infoPanel.message(fmt.Sprintf("New task [%d] created.", task.Id))
		ac.navigateToHomeFunc(false)()
	}
}

func (ac *appController) deleteTaskFunc(taskIDInput *tview.InputField) func() {
	return func() {
		id, err := extractIDFromInput(taskIDInput)
		if err != nil {
			ac.infoPanel.warning(fmt.Sprint("Your input is invalid: ", err))
			return
		}

		t := database.Task{Id: int64(id)}
		err = t.Delete()
		if err != nil {
			ac.infoPanel.error(err)
			return
		}
		ac.infoPanel.message(fmt.Sprintf("Task [%d] deleted.", id))
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
