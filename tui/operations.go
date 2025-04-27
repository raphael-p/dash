package tui

import (
	"fmt"
	"strconv"

	"github.com/raphael-p/datashard/database"
	"github.com/rivo/tview"
)

func (vc viewController) refreshTasks() {
	vc.displayPanel.Clear()

	tasks, err := database.GetTasks("")
	if err != nil {
		vc.logPanel.SetText(fmt.Sprint("[red]", err))
		return
	}

	if len(tasks) == 0 {
		fmt.Fprintln(vc.displayPanel, "No tasks. Good job :)")
		return
	}
	for i, task := range tasks {
		fmt.Fprintf(vc.displayPanel, "[%d] %s\n", i+1, task.Stringify())
	}
}

func (vc viewController) addTask(taskNameInput, taskDescriptionInput *tview.InputField) func() {
	return func() {
		name := taskNameInput.GetText()
		description := taskDescriptionInput.GetText()

		if name == "" || description == "" {
			vc.logPanel.SetText("[red]Please provide a name and description of the task.")
			return
		}

		task, err := database.CreateTask(name, description)
		if err != nil {
			vc.logPanel.SetText(fmt.Sprint("[red]", err))
			return
		}
		vc.logPanel.SetText(fmt.Sprintf("New task [%d] created.", task.Id))
		vc.refreshTasks()
		vc.navigateToHome()()
	}
}

func (vc viewController) deleteTask(taskIDInput *tview.InputField) func() {
	return func() {
		vc.logPanel.SetText("")
		idStr := taskIDInput.GetText()
		if idStr == "" {
			vc.logPanel.SetText("[red]Please enter a task ID to delete.")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			vc.logPanel.SetText("[red]Invalid task ID.")
			return
		}

		t := database.Task{Id: int64(id)}
		err = t.Delete()
		if err != nil {
			vc.logPanel.SetText(fmt.Sprint("[red]", err))
			return
		}
		vc.logPanel.SetText(fmt.Sprintf("Task [%d] deleted.", id))
		vc.refreshTasks()
		vc.navigateToHome()()
	}
}
