package taskform

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TaskPriority int

const (
	PRIORITY_LAST TaskPriority = iota
	PRIORITY_NEXT
	PRIORITY_FIRST
)
const DEFAULT_PRIORITY = PRIORITY_LAST

type FieldSelection struct {
	TaskID, TaskName, TaskDescription, TaskPriority bool
}

type TaskForm struct {
	*tview.Form
	ClearInputs        func()
	GetTaskID          func() string
	GetTaskName        func() string
	GetTaskDescription func() string
	GetTaskPriority    func() TaskPriority
}

func New() *TaskForm {
	form := TaskForm{tview.NewForm(), func() {}, nil, nil, nil, nil}
	return &form
}

func (tf *TaskForm) addInputField(name string) *tview.InputField {
	inputField := tview.NewInputField().SetLabel(name + ": ")
	clearPrevious := tf.ClearInputs
	tf.ClearInputs = func() {
		clearPrevious()
		inputField.SetText("")
	}
	tf.AddFormItem(inputField)
	return inputField
}

func (tf *TaskForm) PromptId() *TaskForm {
	tf.GetTaskID = tf.addInputField("Task ID").GetText
	return tf
}

func (tf *TaskForm) PromptDescription() *TaskForm {
	tf.GetTaskDescription = tf.addInputField("Task Description").GetText
	return tf
}

func (tf *TaskForm) PromptName() *TaskForm {
	tf.GetTaskName = tf.addInputField("Task Name").GetText
	return tf
}

func (tf *TaskForm) PromptPriority() *TaskForm {
	taskPriority := DEFAULT_PRIORITY
	tf.AddDropDown("Priority", []string{"last", "next", "first"}, int(taskPriority), func(_ string, index int) {
		taskPriority = TaskPriority(index)
	})
	tf.GetTaskPriority = func() TaskPriority { return taskPriority }
	return tf
}

func (tf *TaskForm) AddBackButton(back func()) *TaskForm {
	tf.AddButton("Back", back)
	return tf
}

func (tf *TaskForm) AddQuitButton(quit func()) *TaskForm {
	tf.AddButton("Quit", quit)
	return tf
}

func (tf *TaskForm) OnEnter(callback func()) *TaskForm {
	tf.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// used to check that the selected item is not a button
		_, buttonIndex := tf.GetFocusedItemIndex()

		if event.Key() == tcell.KeyEnter && buttonIndex == -1 {
			callback()
			return nil
		}
		return event
	})
	return tf
}

func (tf *TaskForm) OnSubmit(callback func()) *TaskForm {
	tf.AddButton("Submit", callback)
	return tf
}
