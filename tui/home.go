package tui

import (
	"github.com/raphael-p/datashard/logger"
	"github.com/rivo/tview"
)

func Home() {
	app := tview.NewApplication()

	taskList := tview.NewTextView().
		SetWordWrap(true)
	taskList.SetBorder(true).SetTitle("Tasks")

	taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
	taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
	taskIDInput := tview.NewInputField().SetLabel("Task ID to delete: ")

	displayPanel := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
	displayPanel.SetBorder(true)

	form := tview.NewForm().
		AddFormItem(taskNameInput).
		AddFormItem(taskDescriptionInput).
		AddFormItem(taskIDInput).
		AddButton("Add Task", addTask(taskNameInput, taskDescriptionInput, displayPanel, taskList)).
		AddButton("Delete Task", deleteTask(taskIDInput, displayPanel, taskList)).
		AddButton("Exit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Task Management")

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 2, true).
		AddItem(displayPanel, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(taskList, 0, 2, false).
		AddItem(rightFlex, 0, 1, true)

	refreshTasks(displayPanel, taskList)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		logger.Fatal(err.Error())
	}
}
