package controller

import (
	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/taskform"
)

func (c *Controller) openTaskForm(back func(), quit func()) {
	c.infoPanel.Clear()
	fieldSelection := taskform.FieldSelection{TaskID: true}
	openTaskForm := taskform.New(fieldSelection, back, quit)

	openTaskForm.OnEnter(func() {
		c.openTask(openTaskForm.GetTaskID(), back, quit)
		openTaskForm.ClearInputs()
	})

	c.inputPanel.Set("Open Task", openTaskForm)
}

func (c *Controller) addTaskForm(back func(), quit func()) {
	c.infoPanel.Clear()
	fieldSelection := taskform.FieldSelection{TaskName: true, TaskDescription: true}
	addTaskForm := taskform.New(fieldSelection, back, quit)

	addTaskForm.OnEnter(func() {
		ok := c.addTask(addTaskForm.GetTaskName(), addTaskForm.GetTaskDescription())
		addTaskForm.ClearInputs()
		addTaskForm.SetFocus(0)     // manually set focus to first field
		c.app.SetFocus(addTaskForm) // return automatic focus management to form
		if ok {
			back()
		}
	})

	c.inputPanel.Set("Add New Task", addTaskForm)
}

func (c *Controller) removeTaskForm(back func(), quit func()) {
	c.infoPanel.Clear()
	fieldSelection := taskform.FieldSelection{TaskID: true}
	removeTaskForm := taskform.New(fieldSelection, back, quit)

	removeTaskForm.OnEnter(func() {
		c.removeTask(removeTaskForm.GetTaskID())
		removeTaskForm.ClearInputs()
	})

	c.inputPanel.Set("Remove Task", removeTaskForm)
}

func (c *Controller) editTaskForm(task *database.Task, back func(), quit func()) {
	c.infoPanel.Clear()
	fieldSelection := taskform.FieldSelection{TaskName: true, TaskDescription: true}
	editTaskForm := taskform.New(fieldSelection, back, quit)

	editTaskForm.OnEnter(func() {
		ok := c.editTask(task, editTaskForm.GetTaskName(), editTaskForm.GetTaskDescription())
		if ok {
			back()
		}
	})

	c.inputPanel.Set("Edit Task", editTaskForm)
}

func (c *Controller) bumpTaskForm(back func(), quit func()) {
	c.infoPanel.Clear()
	fieldSelection := taskform.FieldSelection{TaskID: true}
	bumpTaskForm := taskform.New(fieldSelection, back, quit)

	bumpTaskForm.OnEnter(func() {
		c.bumpTask(bumpTaskForm.GetTaskID())
		bumpTaskForm.ClearInputs()
	})

	c.inputPanel.Set("Bump Task Priority", bumpTaskForm)
}
