package tui

import (
	"fmt"
	"strconv"

	"github.com/raphael-p/datashard/database"
	"github.com/rivo/tview"
)

func refreshTasks(displayPanel, taskList *tview.TextView) {
	taskList.Clear()

	tasks, err := database.GetTasks("")
	if err != nil {
		displayPanel.SetText(fmt.Sprint("[red]", err))
		return
	}

	if len(tasks) == 0 {
		fmt.Fprintln(taskList, "No tasks. Good job :)")
		return
	}
	for i, task := range tasks {
		fmt.Fprintf(taskList, "[%d] %s\n", i+1, task.Stringify())
	}
}

func addTask(taskNameInput, taskDescriptionInput *tview.InputField, displayPanel, taskList *tview.TextView) func() {
	return func() {
		name := taskNameInput.GetText()
		description := taskDescriptionInput.GetText()

		if name == "" || description == "" {
			displayPanel.SetText("[red]Please provide a name and description of the task.")
			return
		}

		task, err := database.CreateTask(name, description)
		if err != nil {
			displayPanel.SetText(fmt.Sprint("[red]", err))
			return
		}
		displayPanel.SetText(fmt.Sprintf("New task [%d] creatd.", task.Id))

		refreshTasks(displayPanel, taskList)
		taskNameInput.SetText("")
		taskDescriptionInput.SetText("")
	}
}

func deleteTask(taskIDInput *tview.InputField, displayPanel, taskList *tview.TextView) func() {
	return func() {
		displayPanel.SetText("")
		idStr := taskIDInput.GetText()
		if idStr == "" {
			displayPanel.SetText("[red]Please enter a task ID to delete.")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			displayPanel.SetText("[red]Invalid task ID.")
			return
		}

		t := database.Task{Id: int64(id)}
		err = t.Delete()
		if err != nil {
			displayPanel.SetText(fmt.Sprint("[red]", err))
			return
		}
		displayPanel.SetText(fmt.Sprintf("Task [%d] deleted.", id))
		refreshTasks(displayPanel, taskList)
		taskIDInput.SetText("")
	}
}
