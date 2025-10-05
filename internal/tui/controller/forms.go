package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/database"
	"github.com/rivo/tview"
)

func (c *Controller) openTaskForm() {
	openTaskForm := createForm(c)

	taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
	openTaskForm.AddFormItem(taskIDInput)

	onFormEnter(openTaskForm, func() { c.submitOpenTaskInput(taskIDInput) })
	addExitButtonsToForm(c, openTaskForm)

	c.inputPanel.Set("View Task", openTaskForm)
}

func (c *Controller) addTaskForm() {
	addTaskForm := createForm(c)

	taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
	addTaskForm.AddFormItem(taskNameInput)

	taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
	addTaskForm.AddFormItem(taskDescriptionInput)

	onFormEnter(addTaskForm, func() {
		c.submitAddTask(taskNameInput, taskDescriptionInput)
		addTaskForm.SetFocus(0)     // manually set focus to first field
		c.app.SetFocus(addTaskForm) // return automatic focus management to form
	})
	addExitButtonsToForm(c, addTaskForm)

	c.inputPanel.Set("Add New Task", addTaskForm)
}

func (c *Controller) removeTaskForm() {
	removeTaskForm := createForm(c)

	taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
	removeTaskForm.AddFormItem(taskIDInput)

	onFormEnter(removeTaskForm, func() { c.submitRemoveTask(taskIDInput) })
	addExitButtonsToForm(c, removeTaskForm)

	c.inputPanel.Set("Delete Task", removeTaskForm)
}

func (c *Controller) editTaskForm(task database.Task) {
	editTaskForm := createForm(c)

	taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
	editTaskForm.AddFormItem(taskNameInput)

	taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
	editTaskForm.AddFormItem(taskDescriptionInput)

	onFormEnter(editTaskForm, func() {
		c.submitEditTask(
			task,
			taskNameInput,
			taskDescriptionInput,
		)
	})
	addExitButtonsToForm(c, editTaskForm)

	c.inputPanel.Set("Edit Task", editTaskForm)
}

func createForm(c *Controller) *tview.Form {
	c.infoPanel.Clear()
	return tview.NewForm()
}

func addExitButtonsToForm(c *Controller, form *tview.Form) {
	form.AddButton("Back", func() { c.MainMenu() })
	form.AddButton("Quit", func() { c.app.Stop() })
}

func onFormEnter(form *tview.Form, callback func()) {
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// used to check that the selected item is not a button
		_, buttonIndex := form.GetFocusedItemIndex()

		if event.Key() == tcell.KeyEnter && buttonIndex == -1 {
			callback()
			return nil
		}
		return event
	})
}
