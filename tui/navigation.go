package tui

import "github.com/rivo/tview"

func (ac *appController) navigateToHomeFunc(clearInfo bool) func() {
	return func() {
		if clearInfo {
			ac.infoPanel.clear()
		}
		err := ac.displayPanel.listTasks()
		if err != nil {
			ac.infoPanel.error(err)
		}
		navigationList := tview.NewList()

		navigationList.AddItem("View Task", "", '0', ac.navigateToViewTaskFunc())
		navigationList.AddItem("Add Task", "", '1', ac.navigateToAddTaskFunc())
		navigationList.AddItem("Delete Task", "", '2', ac.navigateToDeleteTaskFunc())
		navigationList.AddItem("Quit", "", 'q', func() { ac.app.Stop() })

		ac.inputPanel.set("Manage Your Tasks", navigationList)
	}
}

func (ac *appController) navigateToViewTaskFunc() func() {
	return func() {
		viewTaskForm := tview.NewForm()

		taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
		viewTaskForm.AddFormItem(taskIDInput)
		viewTaskForm.AddButton("View", ac.viewTaskFunc(taskIDInput))
		viewTaskForm.AddButton("Back", ac.navigateToHomeFunc(true))
		viewTaskForm.AddButton("Quit", func() { ac.app.Stop() })

		ac.inputPanel.set("View Task", viewTaskForm)
	}
}

func (ac *appController) navigateToAddTaskFunc() func() {
	return func() {
		addTaskForm := tview.NewForm()

		taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
		addTaskForm.AddFormItem(taskNameInput)

		taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
		addTaskForm.AddFormItem(taskDescriptionInput)

		addTaskForm.AddButton("Add Task", ac.addTaskFunc(taskNameInput, taskDescriptionInput))
		addTaskForm.AddButton("Back", ac.navigateToHomeFunc(true))
		addTaskForm.AddButton("Quit", func() { ac.app.Stop() })

		ac.inputPanel.set("Add New Task", addTaskForm)
	}
}

func (ac *appController) navigateToDeleteTaskFunc() func() {
	return func() {
		deleteTaskForm := tview.NewForm()

		taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
		deleteTaskForm.AddFormItem(taskIDInput)
		deleteTaskForm.AddButton("Delete", ac.deleteTaskFunc(taskIDInput))
		deleteTaskForm.AddButton("Back", ac.navigateToHomeFunc(true))
		deleteTaskForm.AddButton("Quit", func() { ac.app.Stop() })

		ac.inputPanel.set("Delete Task", deleteTaskForm)
	}
}
