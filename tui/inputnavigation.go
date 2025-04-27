package tui

import "github.com/rivo/tview"

func createInputPanel(app *tview.Application, logPanel, displayPanel *tview.TextView) viewController {
	vc := viewController{app, tview.NewForm(), logPanel, displayPanel}
	vc.inputPanel.SetBorder(true)
	vc.navigateToHome()()
	return vc
}

func (vc viewController) addBackButton() {
	vc.inputPanel.AddButton("Back", vc.navigateToHome())
}

func (vc viewController) addExitButton() {
	vc.inputPanel.AddButton("Exit", func() {
		vc.app.Stop()
	})
}

func (vc viewController) navigate(newTitle string, callback func()) func() {
	return func() {
		vc.inputPanel.Clear(true)
		vc.inputPanel.SetTitle(newTitle)
		callback()
		vc.addExitButton()
		vc.app.SetFocus(vc.inputPanel)
	}
}

func (vc viewController) navigateToHome() func() {
	return vc.navigate("Task Management", func() {
		vc.inputPanel.AddButton("Add Task", vc.navigateToAddTask())
		vc.inputPanel.AddButton("Delete Task", vc.navigateToDeleteTask())
	})
}

func (vc viewController) navigateToAddTask() func() {
	return vc.navigate("Add New Task", func() {
		taskNameInput := tview.NewInputField().SetLabel("Task Name: ")
		vc.inputPanel.AddFormItem(taskNameInput)

		taskDescriptionInput := tview.NewInputField().SetLabel("Description: ")
		vc.inputPanel.AddFormItem(taskDescriptionInput)

		vc.inputPanel.AddButton("Add Task", vc.addTask(taskNameInput, taskDescriptionInput))
		vc.addBackButton()
	})
}

func (vc viewController) navigateToDeleteTask() func() {
	return vc.navigate("Delete Task", func() {
		taskIDInput := tview.NewInputField().SetLabel("Task ID: ")
		vc.inputPanel.AddFormItem(taskIDInput)
		vc.inputPanel.AddButton("Delete", vc.deleteTask(taskIDInput))
		vc.addBackButton()
	})
}
