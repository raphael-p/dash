package tui

import (
	"fmt"
	"strconv"

	"github.com/raphael-p/datashard/database"
	"github.com/rivo/tview"
)

func (ac *appController) refreshTasks() {
	ac.displayPanel.Clear()

	tasks, err := database.GetTasks("")
	if err != nil {
		ac.logPanel.SetText(fmt.Sprint("[red]", err))
		return
	}

	if len(tasks) == 0 {
		fmt.Fprintln(ac.displayPanel, "No tasks. Good job :)")
		return
	}
	for i, task := range tasks {
		fmt.Fprintf(ac.displayPanel, "[%d] %s\n", i+1, task.Stringify())
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
		ac.refreshTasks()
		ac.navigateToHomeFunc()()
	}
}

func (ac *appController) deleteTaskFunc(taskIDInput *tview.InputField) func() {
	return func() {
		ac.logPanel.SetText("")
		idStr := taskIDInput.GetText()
		if idStr == "" {
			ac.logPanel.SetText("[red]Please enter a task ID to delete.")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			ac.logPanel.SetText("[red]Invalid task ID.")
			return
		}

		t := database.Task{Id: int64(id)}
		err = t.Delete()
		if err != nil {
			ac.logPanel.SetText(fmt.Sprint("[red]", err))
			return
		}
		ac.logPanel.SetText(fmt.Sprintf("Task [%d] deleted.", id))
		ac.refreshTasks()
		ac.navigateToHomeFunc()()
	}
}
