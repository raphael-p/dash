package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SetupablePanel interface {
	SetTitle(string) *tview.Box
	SetBorder(bool) *tview.Box
	SetBorderColor(tcell.Color) *tview.Box
}

func setupPanel[T SetupablePanel](panel T, name string) {
	panel.SetTitle(name)
	panel.SetBorder(true).SetBorderColor(tcell.ColorLimeGreen)
}

func (vc *viewController) navigateToHome() func() {
	return func() {
		navigationList := tview.NewList()
		setupPanel(navigationList, "Control Panel")

		navigationList.AddItem("Add Task", "", '1', vc.navigateToAddTask())
		navigationList.AddItem("Delete Task", "", '2', vc.navigateToDeleteTask())
		navigationList.AddItem("Quit", "", 'q', func() { vc.app.Stop() })

		vc.switchFocusPanel(navigationList)
	}
}

func (vc *viewController) navigateToAddTask() func() {
	return func() {
		addTaskForm := tview.NewForm()
		setupPanel(addTaskForm, "Add New Task")

		taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
		addTaskForm.AddFormItem(taskNameInput)

		taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
		addTaskForm.AddFormItem(taskDescriptionInput)

		addTaskForm.AddButton("Add Task", vc.addTask(taskNameInput, taskDescriptionInput))
		addTaskForm.AddButton("Back", vc.navigateToHome())
		addTaskForm.AddButton("Quit", func() { vc.app.Stop() })

		vc.switchFocusPanel(addTaskForm)
	}
}

func (vc *viewController) navigateToDeleteTask() func() {
	return func() {
		deleteTaskForm := tview.NewForm()
		setupPanel(deleteTaskForm, "Delete Task")

		taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
		deleteTaskForm.AddFormItem(taskIDInput)
		deleteTaskForm.AddButton("Delete", vc.deleteTask(taskIDInput))
		deleteTaskForm.AddButton("Back", vc.navigateToHome())
		deleteTaskForm.AddButton("Quit", func() { vc.app.Stop() })

		vc.switchFocusPanel(deleteTaskForm)
	}
}
