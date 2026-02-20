package controller

import (
	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/taskform"
)

func (c *Controller) openTaskForm(back func()) {
	c.infoPanel.Clear()

	openTaskForm := taskform.New().PromptId().AddBackButton(back).AddQuitButton(c.app.Stop)
	openTaskForm.OnEnter(func() {
		c.openTask(openTaskForm.GetTaskID(), back)
		openTaskForm.ClearInputs()
	})

	c.inputPanel.Set("Open Task", openTaskForm)
}

func (c *Controller) addTaskForm(back func()) {
	c.infoPanel.Clear()

	addTaskForm := taskform.New().PromptName().PromptDescription().PromptPriority()
	addTaskForm.OnSubmit(func() {
		ok := c.addTask(addTaskForm.GetTaskName(), addTaskForm.GetTaskDescription())
		addTaskForm.ClearInputs()
		addTaskForm.SetFocus(0)     // manually set focus to first field
		c.app.SetFocus(addTaskForm) // return automatic focus management to form
		if ok {
			back()
		}
	}).AddBackButton(back).AddQuitButton(c.app.Stop)

	c.inputPanel.Set("Add New Task", addTaskForm)
}

func (c *Controller) removeTaskForm(back func()) {
	c.infoPanel.Clear()

	removeTaskForm := taskform.New().PromptId().AddBackButton(back).AddQuitButton(c.app.Stop)
	removeTaskForm.OnEnter(func() {
		c.removeTask(removeTaskForm.GetTaskID())
		removeTaskForm.ClearInputs()
	})

	c.inputPanel.Set("Remove Task", removeTaskForm)
}

func (c *Controller) editTaskForm(task *database.Task, back func()) {
	c.infoPanel.Clear()

	editTaskForm := taskform.New().PromptName().PromptDescription().AddBackButton(back).AddQuitButton(c.app.Stop)
	editTaskForm.OnEnter(func() {
		ok := c.editTask(task, editTaskForm.GetTaskName(), editTaskForm.GetTaskDescription())
		if ok {
			back()
		}
	})

	c.inputPanel.Set("Edit Task", editTaskForm)
}

func (c *Controller) bumpTaskForm(back func()) {
	c.infoPanel.Clear()

	bumpTaskForm := taskform.New().PromptId().AddBackButton(back).AddQuitButton(c.app.Stop)
	bumpTaskForm.OnEnter(func() {
		c.bumpTask(bumpTaskForm.GetTaskID())
		bumpTaskForm.ClearInputs()
	})

	c.inputPanel.Set("Bump Task Priority", bumpTaskForm)
}
