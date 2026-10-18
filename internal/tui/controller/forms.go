package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
	"github.com/raphael-p/dash/pkg/tviewcomponents/taskform"
	"github.com/rivo/tview"
)

func (c *Controller) addTaskForm(back func()) {
	c.infoPanel.Clear()

	addTaskForm := taskform.New().
		PromptName(c.config.nameCharLimit).
		PromptDescription(c.config.descriptionCharLimit).
		PromptPriority()
	addTaskForm.OnSubmit(func() {
		ok := c.addTask(addTaskForm.GetTaskName(), addTaskForm.GetTaskDescription(), addTaskForm.GetTaskPriority())
		addTaskForm.SetFocus(0)     // manually set focus to first field
		c.app.SetFocus(addTaskForm) // return automatic focus management to form
		if ok {
			addTaskForm.ClearInputs()
			back()
		}
	}).AddBackButton(back).AddQuitButton(c.app.Stop)

	c.inputPanel.Set("Add New Task", addTaskForm)
}

func (c *Controller) removeTaskForm(back func()) {
	c.infoPanel.Clear()

	removeTaskForm := taskform.New().PromptID().AddBackButton(back).AddQuitButton(c.app.Stop)
	removeTaskForm.OnEnter(func() {
		c.removeTask(removeTaskForm.GetTaskID())
		removeTaskForm.ClearInputs()
	})

	c.inputPanel.Set("Remove Task", removeTaskForm)
}

func (c *Controller) editTaskForm(task *database.Task, back func()) {
	c.infoPanel.Clear()

	editTaskForm := taskform.New().
		PromptName(c.config.nameCharLimit).
		PromptDescription(c.config.descriptionCharLimit).
		AddBackButton(back).
		AddQuitButton(c.app.Stop)
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

	bumpTaskForm := taskform.New().PromptID().AddBackButton(back).AddQuitButton(c.app.Stop)
	bumpTaskForm.OnEnter(func() {
		c.bumpTask(bumpTaskForm.GetTaskID())
		bumpTaskForm.ClearInputs()
	})

	c.inputPanel.Set("Bump Task Priority", bumpTaskForm)
}

func (c *Controller) emptyForm() {
	c.infoPanel.Clear()

	c.inputPanel.GetPanel()

	empty := tview.NewTextView().SetDynamicColors(true)
	empty.SetBorderPadding(1, 1, 2, 2)

	handler := &keybindHandler{c}
	menu := keybindmenu.New().
		SetHighlighColour(tcell.ColorLimeGreen).
		SetFallback(handler.fallback)

	menu.AddKeybind('b', "back", keybindmenu.DefaultBind, handler.backToHome).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit).
		Apply(empty)

	c.inputPanel.Set("Welcome to Dash!", empty)
}
