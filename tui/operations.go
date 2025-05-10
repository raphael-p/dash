package tui

import (
	"fmt"
	"strconv"

	"github.com/raphael-p/datashard/database"
	"github.com/rivo/tview"
)

func (ac *appController) refreshTasks() {
	ac.displayPanel.Clear()
	ac.displayPanel.SetTitle("Tasks")

	tasks, err := database.GetTasks("")
	if err != nil {
		ac.logPanel.SetText(fmt.Sprint("[red]", err))
		return
	}

	if len(tasks) == 0 {
		fmt.Fprintln(ac.displayPanel, "No tasks. Good job :)")
		return
	}
	for _, task := range tasks {
		fmt.Fprintf(ac.displayPanel, "[%d] %s\n", task.Id, task.Name)
	}
}

func (ac *appController) viewTaskFunc(taskIDInput *tview.InputField) func() {
	return func() {
		id, err := extractIDFromInput(taskIDInput)
		if err != nil {
			ac.logPanel.SetText(fmt.Sprint("[red]Your input is invalid: ", err))
			return
		}

		task, err := database.GetTask(int64(id))
		if err != nil {
			ac.logPanel.SetText(fmt.Sprint("[red]Could not retrieve task: ", err))
			return
		}

		ac.displayPanel.Clear()
		ac.displayPanel.SetTitle(fmt.Sprintf("Task %d: %s", task.Id, task.Name))
		task.Display(ac.displayPanel)
	}
}

func (ac *appController) addTaskFunc(taskNameInput, taskDescriptionInput *tview.InputField) func() {
	return func() {
		name := taskNameInput.GetText()
		description := taskDescriptionInput.GetText()

		if name == "" || description == "" {
			ac.logPanel.SetText("[red]Please provide a name and description of the task.")
			return
		}

		task, err := database.CreateTask(name, description)
		if err != nil {
			ac.logPanel.SetText(fmt.Sprint("[red]", err))
			return
		}
		ac.logPanel.SetText(fmt.Sprintf("New task [%d] created.", task.Id))
		ac.navigateToHomeFunc()()
	}
}

func (ac *appController) deleteTaskFunc(taskIDInput *tview.InputField) func() {
	return func() {
		id, err := extractIDFromInput(taskIDInput)
		if err != nil {
			ac.logPanel.SetText(fmt.Sprint("[red]Your input is invalid: ", err))
			return
		}

		t := database.Task{Id: int64(id)}
		err = t.Delete()
		if err != nil {
			ac.logPanel.SetText(fmt.Sprint("[red]", err))
			return
		}
		ac.logPanel.SetText(fmt.Sprintf("Task [%d] deleted.", id))
		ac.navigateToHomeFunc()()
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
