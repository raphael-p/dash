package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (c *Controller) viewTaskForm() {
	viewTaskForm := tview.NewForm()

	taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
	viewTaskForm.AddFormItem(taskIDInput)

	onFormEnter(viewTaskForm, func() { c.submitViewTask(taskIDInput) })
	addExitButtonsToForm(c, viewTaskForm)

	c.inputPanel.Set("View Task", viewTaskForm)
}

func (c *Controller) addTaskForm() {
	addTaskForm := tview.NewForm()

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

func (c *Controller) deleteTaskForm() {
	deleteTaskForm := tview.NewForm()

	taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
	deleteTaskForm.AddFormItem(taskIDInput)

	onFormEnter(deleteTaskForm, func() { c.submitDeleteTask(taskIDInput) })
	addExitButtonsToForm(c, deleteTaskForm)

	c.inputPanel.Set("Delete Task", deleteTaskForm)
}

func addExitButtonsToForm(c *Controller, form *tview.Form) {
	form.AddButton("Back", func() { c.infoPanel.Clear(); c.MainMenu() })
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
