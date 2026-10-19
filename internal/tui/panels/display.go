package panels

import (
	"fmt"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/stringpad"
	"github.com/rivo/tview"
)

type taskPage struct {
	lastTask database.Task
	hasNext  bool
}

type displayPanelTab struct {
	pages            map[uint]taskPage
	currentPageIndex uint
}

func newDisplayPanelTab() *displayPanelTab {
	return &displayPanelTab{map[uint]taskPage{}, 0}
}

type DisplayPanel struct {
	textView     *tview.TextView
	activeTab    *displayPanelTab
	completedTab *displayPanelTab
	mode         database.TaskMode
}

func NewDisplayPanel() *DisplayPanel {
	panel := tview.NewTextView().SetWordWrap(true)
	panel.SetBorder(true).SetBorderPadding(1, 1, 2, 2)
	return &DisplayPanel{panel, newDisplayPanelTab(), newDisplayPanelTab(), database.ActiveTasks}
}

func (dp *DisplayPanel) GetPrimitive() tview.Primitive {
	return dp.textView
}

func (dp *DisplayPanel) ShowCompleted() bool {
	return dp.mode == database.CompletedTasks
}

func (dp *DisplayPanel) getTab() *displayPanelTab {
	if dp.mode == database.CompletedTasks {
		return dp.completedTab
	}
	return dp.activeTab
}

func (dp *DisplayPanel) ToggleTaskMode() {
	if dp.mode == database.CompletedTasks {
		dp.mode = database.ActiveTasks
	} else {
		dp.mode = database.CompletedTasks
	}
}

func (dp *DisplayPanel) ResetPagination() {
	dp.getTab().currentPageIndex = 0
}

func (dp *DisplayPanel) GetCurrentPage() error {
	return dp.listTasks(dp.getTab().currentPageIndex)
}

func (dp *DisplayPanel) getPreviousPage() error {
	tab := dp.getTab()
	if tab.currentPageIndex == 0 {
		return nil // intentional noop
	}

	pageIndex := tab.currentPageIndex - 1
	err := dp.listTasks(pageIndex)
	if err == nil {
		tab.currentPageIndex = pageIndex
		dp.textView.ScrollToEnd()
	}
	return err
}

func (dp *DisplayPanel) GetNextPage() error {
	tab := dp.getTab()
	currentPage := tab.pages[tab.currentPageIndex]
	if !currentPage.hasNext {
		return nil // intentional noop
	}

	pageIndex := tab.currentPageIndex + 1
	err := dp.listTasks(pageIndex)
	if err == nil {
		tab.currentPageIndex = pageIndex
		dp.textView.ScrollToBeginning()
	}
	return err
}

func (dp *DisplayPanel) listTasks(pageIndex uint) error {
	dp.textView.Clear()
	if dp.mode == database.CompletedTasks {
		dp.textView.SetTitle(fmt.Sprintf(" Completed Tasks (page %d) ", pageIndex+1))
	} else {
		dp.textView.SetTitle(fmt.Sprintf(" Active Tasks (page %d) ", pageIndex+1))
	}

	pages := dp.getTab().pages
	var cursor database.TaskCursor
	if pageIndex != 0 {
		cursor.ID = pages[pageIndex-1].lastTask.ID
		if dp.mode == database.CompletedTasks {
			cursor.Timestamp = pages[pageIndex-1].lastTask.CompletedAt
		} else {
			cursor.Timestamp = pages[pageIndex-1].lastTask.PriorityBumpedAt
		}

	}

	tasks, hasNext, err := database.GetTasksPaginated(dp.mode, cursor)
	if err != nil {
		return fmt.Errorf("could not retrieve tasks: %s", err)
	}

	if len(tasks) == 0 {
		fmt.Fprintln(dp.textView, "No tasks. Good job :)")
		return nil
	}

	// store pagination metadata
	pages[pageIndex] = taskPage{tasks[len(tasks)-1], hasNext}

	if pageIndex > 0 {
		fmt.Fprintf(dp.textView, "\n↑ show previous page\n\n")
	}
	for idx, task := range tasks {
		fmt.Fprint(
			dp.textView,
			stringpad.RightPad(fmt.Sprintf("[%d]", task.ID), 8),
			task.Name,
		)
		if idx+1 < len(tasks) {
			fmt.Fprint(dp.textView, "\n")
		}
	}
	if hasNext {
		fmt.Fprintf(dp.textView, "\n\n↓ show next page\n")
	}
	return nil
}

func (dp *DisplayPanel) showTask(task database.Task) (database.Task, error) {
	dp.textView.Clear()
	dp.textView.SetTitle(fmt.Sprintf("Task %d: %s", task.ID, task.Name))
	if task.Description != "" {
		fmt.Fprintf(dp.textView, "%s\n\n", task.Description)
	}

	createdAt, updatedAt, completedAt := task.GetFormattedTimes()
	if completedAt != "" {
		fmt.Fprintf(dp.textView, "Completed at: %s\n", completedAt)
	}
	fmt.Fprintf(dp.textView, "Created at: %s\n", createdAt)
	fmt.Fprintf(dp.textView, "Last updated: %s\n", updatedAt)
	return task, nil
}

func (dp *DisplayPanel) ShowTopTask() (database.Task, error) {
	task, err := database.GetTopTask()
	if err != nil {
		return database.Task{}, fmt.Errorf("could not retrieve top task: %s", err)
	}

	return dp.showTask(task)
}

func (dp *DisplayPanel) ShowTaskByID(taskID int64) (database.Task, error) {
	task, err := database.GetTask(int64(taskID))
	if err != nil {
		return database.Task{}, fmt.Errorf("could not retrieve task [%d]: %s", taskID, err)
	}

	return dp.showTask(task)
}

func (dp *DisplayPanel) ScrollDown() error {
	offset, _ := dp.textView.GetScrollOffset()

	// determine whether to fetch the next page
	lineCount := dp.textView.GetWrappedLineCount()
	_, _, _, height := dp.textView.GetInnerRect()
	if offset+height >= lineCount {
		return dp.GetNextPage()
	}

	dp.textView.ScrollTo(offset+1, 0)
	return nil
}

func (dp *DisplayPanel) ScrollUp() error {
	offset, _ := dp.textView.GetScrollOffset()

	// determine whether to fetch the previous page
	if offset == 0 {
		return dp.getPreviousPage()
	}

	dp.textView.ScrollTo(offset-1, 0)
	return nil
}
