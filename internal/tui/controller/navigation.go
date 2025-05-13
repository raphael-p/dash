package controller

import "github.com/rivo/tview"

func (ac *Controller) NavigateToHome() {
	err := ac.displayPanel.ListTasks()
	if err != nil {
		ac.infoPanel.Error(err)
	}
	navigationList := tview.NewList()

	navigationList.AddItem("View Task", "", '0', func() { ac.navigateToViewTask() })
	navigationList.AddItem("Add Task", "", '1', func() { ac.navigateToAddTask() })
	navigationList.AddItem("Delete Task", "", '2', func() { ac.navigateToDeleteTask() })
	navigationList.AddItem("Quit", "", 'q', func() { ac.app.Stop() })

	ac.inputPanel.Set("Manage Your Tasks", navigationList)
}

func (ac *Controller) navigateToViewTask() {
	viewTaskForm := tview.NewForm()

	taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
	viewTaskForm.AddFormItem(taskIDInput)
	viewTaskForm.AddButton("View", func() { ac.viewTask(taskIDInput) })
	viewTaskForm.AddButton("Back", func() { ac.infoPanel.Clear(); ac.NavigateToHome() })
	viewTaskForm.AddButton("Quit", func() { ac.app.Stop() })

	ac.inputPanel.Set("View Task", viewTaskForm)
}

func (ac *Controller) navigateToAddTask() {
	addTaskForm := tview.NewForm()

	taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
	addTaskForm.AddFormItem(taskNameInput)

	taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
	addTaskForm.AddFormItem(taskDescriptionInput)

	addTaskForm.AddButton("Add Task", func() { ac.addTask(taskNameInput, taskDescriptionInput) })
	addTaskForm.AddButton("Back", func() { ac.infoPanel.Clear(); ac.NavigateToHome() })
	addTaskForm.AddButton("Quit", func() { ac.app.Stop() })

	ac.inputPanel.Set("Add New Task", addTaskForm)
}

func (ac *Controller) navigateToDeleteTask() {
	deleteTaskForm := tview.NewForm()

	taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
	deleteTaskForm.AddFormItem(taskIDInput)
	deleteTaskForm.AddButton("Delete", func() { ac.deleteTask(taskIDInput) })
	deleteTaskForm.AddButton("Back", func() { ac.infoPanel.Clear(); ac.NavigateToHome() })
	deleteTaskForm.AddButton("Quit", func() { ac.app.Stop() })

	ac.inputPanel.Set("Delete Task", deleteTaskForm)
}
