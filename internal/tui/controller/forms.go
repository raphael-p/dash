package controller

import (
	"github.com/raphael-p/datashard/internal/database"
)

func (c *Controller) openTaskForm() {
	openTaskForm := createTaskForm(c, fieldSelection{taskID: true})
	onFormEnter(openTaskForm, func() {
		c.openTask(openTaskForm.getTaskID())
		openTaskForm.clearInputs()
	})
	c.inputPanel.Set("Open Task", openTaskForm)
}

func (c *Controller) addTaskForm() {
	addTaskForm := createTaskForm(c, fieldSelection{taskName: true, taskDescription: true})
	onFormEnter(addTaskForm, func() {
		c.addTask(addTaskForm.getTaskName(), addTaskForm.getTaskDescription())
		addTaskForm.clearInputs()
		addTaskForm.SetFocus(0)     // manually set focus to first field
		c.app.SetFocus(addTaskForm) // return automatic focus management to form
	})
	c.inputPanel.Set("Add New Task", addTaskForm)
}

func (c *Controller) removeTaskForm() {
	removeTaskForm := createTaskForm(c, fieldSelection{taskID: true})
	onFormEnter(removeTaskForm, func() {
		c.removeTask(removeTaskForm.getTaskID())
		removeTaskForm.clearInputs()
	})
	c.inputPanel.Set("Remove Task", removeTaskForm)
}

func (c *Controller) editTaskForm(task database.Task) {
	editTaskForm := createTaskForm(c, fieldSelection{taskName: true, taskDescription: true})
	onFormEnter(editTaskForm, func() {
		c.editTask(task, editTaskForm.getTaskName(), editTaskForm.getTaskDescription())
	})
	c.inputPanel.Set("Edit Task", editTaskForm)
}

func (c *Controller) bumpTaskForm() {
	bumpTaskForm := createTaskForm(c, fieldSelection{taskID: true})
	onFormEnter(bumpTaskForm, func() {
		c.bumpTask(bumpTaskForm.getTaskID())
		bumpTaskForm.clearInputs()
	})
	c.inputPanel.Set("Bump Task Priority", bumpTaskForm)
}
