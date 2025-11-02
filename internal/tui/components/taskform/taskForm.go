package taskform

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FieldSelection struct {
	TaskID, TaskName, TaskDescription bool
}

type TaskForm struct {
	*tview.Form
	ClearInputs        func()
	GetTaskID          func() string
	GetTaskName        func() string
	GetTaskDescription func() string
}

func New(fieldSelection FieldSelection, back func(), quit func()) *TaskForm {
	form := TaskForm{tview.NewForm(), nil, nil, nil, nil}

	var taskIDInput *tview.InputField
	if fieldSelection.TaskID {
		taskIDInput = tview.NewInputField().SetLabel("Task ID: ")
		form.AddFormItem(taskIDInput)
	}

	var taskNameInput *tview.InputField
	if fieldSelection.TaskName {
		taskNameInput = tview.NewInputField().SetLabel("Task Name: ")
		form.AddFormItem(taskNameInput)
	}

	var taskDescriptionInput *tview.InputField
	if fieldSelection.TaskDescription {
		taskDescriptionInput = tview.NewInputField().SetLabel("Task Description: ")
		form.AddFormItem(taskDescriptionInput)
	}

	form.ClearInputs = func() {
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
	form.AddButton("Back", back)
	form.AddButton("Quit", quit)

	if taskIDInput != nil {
		form.GetTaskID = taskIDInput.GetText
	}
	if taskNameInput != nil {
		form.GetTaskName = taskNameInput.GetText
	}
	if taskDescriptionInput != nil {
		form.GetTaskDescription = taskDescriptionInput.GetText
	}
	return &form
}

func (tf *TaskForm) OnEnter(callback func()) {
	tf.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// used to check that the selected item is not a button
		_, buttonIndex := tf.GetFocusedItemIndex()

		if event.Key() == tcell.KeyEnter && buttonIndex == -1 {
			callback()
			return nil
		}
		return event
	})
}
