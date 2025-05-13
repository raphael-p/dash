package panels

import (
	"fmt"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/rivo/tview"
)

type DisplayPanel struct {
	panel *tview.TextView
}

func NewDisplayPanel() *DisplayPanel {
	panel := tview.NewTextView().SetWordWrap(true)
	panel.SetBorder(true).SetBorderPadding(1, 1, 2, 2)
	return &DisplayPanel{panel}
}

func (dp *DisplayPanel) GetPanel() *tview.TextView {
	return dp.panel
}

func (dp *DisplayPanel) ListTasks() error {
	dp.panel.Clear()
	dp.panel.SetTitle(" Tasks ")

	tasks, err := database.GetTasks("")
	if err != nil {
		return fmt.Errorf("could not retrieve tasks: %s", err)
	}

	if len(tasks) == 0 {
		fmt.Fprintln(dp.panel, "No tasks. Good job :)")
		return nil
	}

	for _, task := range tasks {
		fmt.Fprintf(dp.panel, "[%d] %s\n", task.Id, task.Name)
	}
	return nil
}

func (dp *DisplayPanel) ShowTask(taskId int64) error {
	task, err := database.GetTask(int64(taskId))
	if err != nil {
		return fmt.Errorf("could not retrieve task %d: %s", taskId, err)
	}

	dp.panel.Clear()
	dp.panel.SetTitle(fmt.Sprintf("Task %d: %s", task.Id, task.Name))
	if task.Description != "" {
		fmt.Fprintf(dp.panel, "%s\n\n", task.Description)
	}

	createdAt, updatedAt, completedAt := task.GetFormattedTimes()
	if completedAt != "" {
		fmt.Fprintf(dp.panel, "Completed at: %s\n", completedAt)
	}
	fmt.Fprintf(dp.panel, "Created at: %s\n", createdAt)
	fmt.Fprintf(dp.panel, "Last updated: %s\n", updatedAt)
	return nil
}
