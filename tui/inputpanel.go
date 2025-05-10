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

func (ac *appController) switchInputPanel(newPrimitive tview.Primitive) {
	oldPrimitive := ac.inputPanel
	if oldPrimitive != nil {
		ac.rightFlex.RemoveItem(oldPrimitive)
	}
	ac.rightFlex.Clear().AddItem(newPrimitive, 0, 2, true)
	ac.rightFlex.AddItem(ac.logPanel, 0, 1, false)
	ac.app.SetFocus(newPrimitive)
	ac.inputPanel = newPrimitive
}

func (ac *appController) navigateToHomeFunc() func() {
	return func() {
		ac.refreshTasks()
		navigationList := tview.NewList()
		setupPanel(navigationList, "Manage Your Tasks")

		navigationList.AddItem("View Task", "", '0', ac.navigateToViewTaskFunc())
		navigationList.AddItem("Add Task", "", '1', ac.navigateToAddTaskFunc())
		navigationList.AddItem("Delete Task", "", '2', ac.navigateToDeleteTaskFunc())
		navigationList.AddItem("Quit", "", 'q', func() { ac.app.Stop() })

		ac.switchInputPanel(navigationList)
	}
}

func (ac *appController) navigateToViewTaskFunc() func() {
	return func() {
		viewTaskForm := tview.NewForm()
		setupPanel(viewTaskForm, "View Task")

		taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
		viewTaskForm.AddFormItem(taskIDInput)
		viewTaskForm.AddButton("View", ac.viewTaskFunc(taskIDInput))
		viewTaskForm.AddButton("Back", ac.navigateToHomeFunc())
		viewTaskForm.AddButton("Quit", func() { ac.app.Stop() })

		ac.switchInputPanel(viewTaskForm)
	}
}

func (ac *appController) navigateToAddTaskFunc() func() {
	return func() {
		addTaskForm := tview.NewForm()
		setupPanel(addTaskForm, "Add New Task")

		taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
		addTaskForm.AddFormItem(taskNameInput)

		taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
		addTaskForm.AddFormItem(taskDescriptionInput)

		addTaskForm.AddButton("Add Task", ac.addTaskFunc(taskNameInput, taskDescriptionInput))
		addTaskForm.AddButton("Back", ac.navigateToHomeFunc())
		addTaskForm.AddButton("Quit", func() { ac.app.Stop() })

		ac.switchInputPanel(addTaskForm)
	}
}

func (ac *appController) navigateToDeleteTaskFunc() func() {
	return func() {
		deleteTaskForm := tview.NewForm()
		setupPanel(deleteTaskForm, "Delete Task")

		taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
		deleteTaskForm.AddFormItem(taskIDInput)
		deleteTaskForm.AddButton("Delete", ac.deleteTaskFunc(taskIDInput))
		deleteTaskForm.AddButton("Back", ac.navigateToHomeFunc())
		deleteTaskForm.AddButton("Quit", func() { ac.app.Stop() })

		ac.switchInputPanel(deleteTaskForm)
	}
}
