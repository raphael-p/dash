package taskform

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TaskPriority int

const (
	PriorityLast TaskPriority = iota
	PriorityNext
	PriorityFirst
)
const DefaultPriority = PriorityLast

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

func stopAtCharLimit(textGetter func() string, charLimit int) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		isCharEvent := event.Key() == tcell.KeyRune
		if isCharEvent && len(textGetter()) >= charLimit {
			return nil
		}
		return event
	}
}

func (tf *TaskForm) addInputField(name string, charLimit int) *tview.InputField {
	inputField := tview.NewInputField().SetLabel(name + ": ")
	if charLimit > 0 {
		inputField.SetInputCapture(stopAtCharLimit(inputField.GetText, charLimit))
	}
	clearPrevious := tf.ClearInputs
	tf.ClearInputs = func() {
		clearPrevious()
		inputField.SetText("")
	}
	tf.AddFormItem(inputField)
	return inputField
}

func (tf *TaskForm) PromptID() *TaskForm {
	tf.GetTaskID = tf.addInputField("Task ID", 0).GetText
	return tf
}

func (tf *TaskForm) PromptDescription(charLimit int) *TaskForm {
	tf.GetTaskDescription = tf.addInputField("Task Description", charLimit).GetText
	return tf
}

func (tf *TaskForm) PromptName(charLimit int) *TaskForm {
	tf.GetTaskName = tf.addInputField("Task Name", charLimit).GetText
	return tf
}

func (tf *TaskForm) PromptPriority() *TaskForm {
	taskPriority := DefaultPriority
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
