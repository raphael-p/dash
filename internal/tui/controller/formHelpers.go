package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type fieldSelection struct {
	taskID, taskName, taskDescription bool
}

type taskForm struct {
	*tview.Form
	clearInputs        func()
	getTaskID          func() string
	getTaskName        func() string
	getTaskDescription func() string
}

func createTaskForm(c *Controller, fieldSelection fieldSelection) taskForm {
	c.infoPanel.Clear()
	form := taskForm{tview.NewForm(), nil, nil, nil, nil}

	var taskIDInput *tview.InputField
	if fieldSelection.taskID {
		taskIDInput = tview.NewInputField().SetLabel("Task ID: ")
		form.AddFormItem(taskIDInput)
	}

	var taskNameInput *tview.InputField
	if fieldSelection.taskName {
		taskNameInput = tview.NewInputField().SetLabel("Task Name: ")
		form.AddFormItem(taskNameInput)
	}

	var taskDescriptionInput *tview.InputField
	if fieldSelection.taskDescription {
		taskDescriptionInput = tview.NewInputField().SetLabel("Task Description: ")
		form.AddFormItem(taskDescriptionInput)
	}

	form.clearInputs = func() {
		if taskIDInput != nil {
			taskIDInput.SetText("")
		}
		if taskNameInput != nil {
			taskNameInput.SetText("")
		}
		if taskDescriptionInput != nil {
			taskDescriptionInput.SetText("")
		}
	}
	addExitButtonsToForm(c, form.Form)

	if taskIDInput != nil {
		form.getTaskID = taskIDInput.GetText
	}
	if taskNameInput != nil {
		form.getTaskName = taskNameInput.GetText
	}
	if taskDescriptionInput != nil {
		form.getTaskDescription = taskDescriptionInput.GetText
	}
	return form
}

func addExitButtonsToForm(c *Controller, form *tview.Form) {
	form.AddButton("Back", func() { c.Home() })
	form.AddButton("Quit", func() { c.app.Stop() })
}

func onFormEnter(form taskForm, callback func()) {
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
