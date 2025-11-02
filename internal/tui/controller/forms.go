package controller

import (
	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/internal/tui/components/taskform"
)

func newTaskForm(c *Controller, fieldSelection taskform.FieldSelection) *taskform.TaskForm {
	c.infoPanel.Clear()
	return taskform.New(fieldSelection, func() { c.Home() }, func() { c.app.Stop() })
}

func (c *Controller) openTaskForm() {
	openTaskForm := newTaskForm(c, taskform.FieldSelection{TaskID: true})
	openTaskForm.OnEnter(func() {
		c.openTask(openTaskForm.GetTaskID())
		openTaskForm.ClearInputs()
	})
	c.inputPanel.Set("Open Task", openTaskForm)
}

func (c *Controller) addTaskForm() {
	addTaskForm := newTaskForm(c, taskform.FieldSelection{TaskName: true, TaskDescription: true})
	addTaskForm.OnEnter(func() {
		c.addTask(addTaskForm.GetTaskName(), addTaskForm.GetTaskDescription())
		addTaskForm.ClearInputs()
		addTaskForm.SetFocus(0)     // manually set focus to first field
		c.app.SetFocus(addTaskForm) // return automatic focus management to form
	})
	c.inputPanel.Set("Add New Task", addTaskForm)
}

func (c *Controller) removeTaskForm() {
	removeTaskForm := newTaskForm(c, taskform.FieldSelection{TaskID: true})
	removeTaskForm.OnEnter(func() {
		c.removeTask(removeTaskForm.GetTaskID())
		removeTaskForm.ClearInputs()
	})
	c.inputPanel.Set("Remove Task", removeTaskForm)
}

func (c *Controller) editTaskForm(task database.Task) {
	editTaskForm := newTaskForm(c, taskform.FieldSelection{TaskName: true, TaskDescription: true})
	editTaskForm.OnEnter(func() {
		c.editTask(task, editTaskForm.GetTaskName(), editTaskForm.GetTaskDescription())
	})
	c.inputPanel.Set("Edit Task", editTaskForm)
}

func (c *Controller) bumpTaskForm() {
	bumpTaskForm := newTaskForm(c, taskform.FieldSelection{TaskID: true})
	bumpTaskForm.OnEnter(func() {
		c.bumpTask(bumpTaskForm.GetTaskID())
		bumpTaskForm.ClearInputs()
	})
	c.inputPanel.Set("Bump Task Priority", bumpTaskForm)
}
