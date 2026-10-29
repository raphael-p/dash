package panels

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/internal/tui/components"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/rivo/tview"
)

const minimumTasksPerPage int = 20

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
	frame        *components.Frame
	activeTab    *displayPanelTab
	completedTab *displayPanelTab
	pausedTab    *displayPanelTab
	mode         database.TaskMode
	rowSpace     int
}

func NewDisplayPanel(headerFooterColour tcell.Color) *DisplayPanel {
	frame := components.NewFrame(headerFooterColour)
	return &DisplayPanel{frame, newDisplayPanelTab(), newDisplayPanelTab(), newDisplayPanelTab(), database.ActiveTasks, 0}
}

func (dp *DisplayPanel) GetPrimitive() tview.Primitive {
	return dp.frame
}

func (dp *DisplayPanel) SetPrimitive(primitive tview.Primitive) *DisplayPanel {
	dp.frame.Clear().SetPrimitive(primitive)
	return dp
}

func (dp *DisplayPanel) GetTaskMode() database.TaskMode {
	return dp.mode
}

func (dp *DisplayPanel) HasTasks() bool {
	for _, page := range dp.getTab().pages {
		if page.lastTask.ID != 0 {
			return true
		}
	}
	return false
}

func (dp *DisplayPanel) getTab() *displayPanelTab {
	switch dp.mode {
	case database.CompletedTasks:
		return dp.completedTab
	case database.PausedTasks:
		return dp.pausedTab
	default:
		return dp.activeTab
	}
}

func (dp *DisplayPanel) ToggleTaskMode() {
	dp.mode = dp.mode.GetNext()
}

func (dp *DisplayPanel) ResetPagination() {
	dp.getTab().currentPageIndex = 0
}

func (dp *DisplayPanel) GetCurrentPage() error {
	return dp.listTasks(dp.getTab().currentPageIndex, false)
}

func (dp *DisplayPanel) getPreviousPage() error {
	tab := dp.getTab()
	if tab.currentPageIndex == 0 {
		return nil // intentional noop
	}

	pageIndex := tab.currentPageIndex - 1
	err := dp.listTasks(pageIndex, true)
	if err == nil {
		tab.currentPageIndex = pageIndex
	}
	return err
}

func (dp *DisplayPanel) getNextPage() error {
	tab := dp.getTab()
	currentPage := tab.pages[tab.currentPageIndex]
	if !currentPage.hasNext {
		return nil // intentional noop
	}

	pageIndex := tab.currentPageIndex + 1
	err := dp.listTasks(pageIndex, false)
	if err == nil {
		tab.currentPageIndex = pageIndex
	}
	return err
}

func (dp *DisplayPanel) SetTitle(title string) *DisplayPanel {
	dp.frame.SetTitle(tview.Escape(title))
	return dp
}

func (dp *DisplayPanel) listTasks(pageIndex uint, scrollToEnd bool) error {
	dp.SetTitle(fmt.Sprintf(" %s Tasks (page %d) ", dp.mode.ToString(), pageIndex+1))

	pages := dp.getTab().pages
	var cursor database.TaskCursor
	if pageIndex != 0 {
		cursor.ID = pages[pageIndex-1].lastTask.ID
		switch dp.mode {
		case database.ActiveTasks:
			cursor.Timestamp = pages[pageIndex-1].lastTask.PriorityBumpedAt
		case database.CompletedTasks:
			cursor.Timestamp = pages[pageIndex-1].lastTask.CompletedAt
		case database.PausedTasks:
			cursor.Timestamp = pages[pageIndex-1].lastTask.PausedAt
		}
	} else {
		dp.rowSpace = dp.frame.GetAvailableRows()
	}

	tasks, hasNext, err := database.GetTasksPaginated(dp.mode, cursor, max(minimumTasksPerPage, dp.rowSpace))
	if err != nil {
		return fmt.Errorf("could not retrieve tasks: %w", err)
	}

	if len(tasks) == 0 {
		delete(pages, pageIndex)
		if pageIndex == 0 {
			dp.SetPrimitive(tview.NewTextView().SetText("No tasks. Good job :)"))
			return nil
		}
		return dp.listTasks(pageIndex-1, scrollToEnd)
	}

	// store pagination metadata
	pages[pageIndex] = taskPage{tasks[len(tasks)-1], hasNext}

	table := components.NewTable(tasks)
	dp.SetPrimitive(table)
	if pageIndex > 0 {
		dp.frame.AddHeader("↑ previous page")
	}
	if hasNext {
		dp.frame.AddFooter("↓ next page")
	}
	if scrollToEnd {
		table.ScrollToEnd()
	} else {
		table.ScrollToBeginning()
	}
	return nil
}

func (dp *DisplayPanel) ShowTask(task database.Task) (database.Task, error) {
	dp.SetTitle(fmt.Sprintf("Task %d: %s", task.ID, task.Name))
	textView := tview.NewTextView()
	dp.SetPrimitive(textView)
	if task.Description != "" {
		fmt.Fprintf(textView, "%s\n\n", task.Description)
	}

	dp.frame.AddDateFooter(strutil.FormatTime(task.UpdatedAt), "updated at")
	dp.frame.AddDateFooter(strutil.FormatTime(task.CompletedAt), "completed at")
	dp.frame.AddDateFooter(strutil.FormatTime(task.PausedAt), "paused at")
	return task, nil
}

func (dp *DisplayPanel) ScrollTable(direction components.ScrollDirection) error {
	table, ok := dp.frame.GetPrimitive().(*components.Table)
	if !ok {
		return errors.New("task list is not available")
	}
	if table.Scroll(components.ScrollDirection(direction)) {
		if direction == components.ScrollDown {
			return dp.getNextPage()
		} else {
			return dp.getPreviousPage()
		}
	}
	return nil
}

func (dp *DisplayPanel) ScrollTask(direction components.ScrollDirection) error {
	textView, ok := dp.frame.GetPrimitive().(*tview.TextView)
	if !ok {
		return errors.New("task is not available")
	}
	offset, _ := textView.GetScrollOffset()
	if direction == components.ScrollDown {
		textView.ScrollTo(offset+1, 0)
	} else if offset != 0 {
		textView.ScrollTo(offset-1, 0)
	}
	return nil
}
