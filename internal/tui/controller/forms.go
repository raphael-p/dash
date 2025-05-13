package controller

import "github.com/rivo/tview"

func (c *Controller) viewTaskForm() {
	viewTaskForm := tview.NewForm()

	taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
	viewTaskForm.AddFormItem(taskIDInput)

	viewTaskForm.AddButton("View", func() { c.submitViewTask(taskIDInput) })
	addExitButtonsToForm(c, viewTaskForm)

	c.inputPanel.Set("View Task", viewTaskForm)
}

func (c *Controller) addTaskForm() {
	addTaskForm := tview.NewForm()

	taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
	addTaskForm.AddFormItem(taskNameInput)

	taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
	addTaskForm.AddFormItem(taskDescriptionInput)

	addTaskForm.AddButton("Add Task", func() { c.submitAddTask(taskNameInput, taskDescriptionInput) })
	addExitButtonsToForm(c, addTaskForm)

	c.inputPanel.Set("Add New Task", addTaskForm)
}

func (c *Controller) deleteTaskForm() {
	deleteTaskForm := tview.NewForm()

	taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
	deleteTaskForm.AddFormItem(taskIDInput)

	deleteTaskForm.AddButton("Delete", func() { c.submitDeleteTask(taskIDInput) })
	addExitButtonsToForm(c, deleteTaskForm)

	c.inputPanel.Set("Delete Task", deleteTaskForm)
}

func addExitButtonsToForm(c *Controller, form *tview.Form) {
	form.AddButton("Back", func() { c.infoPanel.Clear(); c.MainMenu() })
	form.AddButton("Quit", func() { c.app.Stop() })
}
