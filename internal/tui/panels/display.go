package panels

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/stringpad"
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
	frame              *tview.Frame
	activeTab          *displayPanelTab
	completedTab       *displayPanelTab
	mode               database.TaskMode
	rowSpace           int
	headerFooterColour tcell.Color
}

func NewDisplayPanel(headerFooterColour tcell.Color) *DisplayPanel {
	frame := newFrame()
	return &DisplayPanel{frame, newDisplayPanelTab(), newDisplayPanelTab(), database.ActiveTasks, 0, headerFooterColour}
}

func (dp *DisplayPanel) GetPrimitive() tview.Primitive {
	return dp.frame
}

func (dp *DisplayPanel) SetPrimitive(primitive tview.Primitive) *DisplayPanel {
	dp.frame.Clear().SetPrimitive(primitive)
	return dp
}

func (dp *DisplayPanel) ShowCompleted() bool {
	return dp.mode == database.CompletedTasks
}

func (dp *DisplayPanel) IsEmpty() bool {
	for _, page := range dp.getTab().pages {
		if page.lastTask.ID != 0 {
			return false
		}
	}
	return true
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

func (dp *DisplayPanel) GetNextPage() error {
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
	dp.frame.SetTitle(title)
	return dp
}

func (dp *DisplayPanel) listTasks(pageIndex uint, scrollToEnd bool) error {
	isCompleted := dp.mode == database.CompletedTasks
	if isCompleted {
		dp.SetTitle(fmt.Sprintf(" Completed Tasks (page %d) ", pageIndex+1))
	} else {
		dp.SetTitle(fmt.Sprintf(" Active Tasks (page %d) ", pageIndex+1))
	}

	pages := dp.getTab().pages
	var cursor database.TaskCursor
	if pageIndex != 0 {
		cursor.ID = pages[pageIndex-1].lastTask.ID
		if isCompleted {
			cursor.Timestamp = pages[pageIndex-1].lastTask.CompletedAt
		} else {
			cursor.Timestamp = pages[pageIndex-1].lastTask.PriorityBumpedAt
		}
	} else {
		dp.rowSpace = getFrameAvailableRows(dp.frame)
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

	table := makeTaskTable(tasks, isCompleted)
	dp.SetPrimitive(table)
	if pageIndex > 0 {
		dp.frame.AddText("↑ previous page", true, tview.AlignLeft, dp.headerFooterColour)
	}
	if hasNext {
		dp.frame.AddText("↓ next page", false, tview.AlignLeft, dp.headerFooterColour)
	}

	if scrollToEnd {
		table.ScrollToEnd()
	} else {
		table.ScrollToBeginning()
	}
	return nil
}

func (dp *DisplayPanel) showTask(task database.Task) (database.Task, error) {
	dp.SetTitle(fmt.Sprintf("Task %d: %s", task.ID, task.Name))
	textView := tview.NewTextView()
	dp.SetPrimitive(textView)
	if task.Description != "" {
		fmt.Fprintf(textView, "%s\n\n", task.Description)
	}

	createdAt, updatedAt, completedAt := task.GetFormattedTimes()

	addDateFooter := func(label, value string) {
		dp.frame.AddText(stringpad.RightPad(label, 14)+value, false, tview.AlignLeft, dp.headerFooterColour)
	}
	if completedAt != "" {
		addDateFooter("completed at", completedAt)
	}
	addDateFooter("created at", createdAt)
	addDateFooter("updated at", updatedAt)
	return task, nil
}

func (dp *DisplayPanel) ShowTopTask() (database.Task, error) {
	task, err := database.GetTopTask()
	if err != nil {
		return database.Task{}, fmt.Errorf("could not retrieve top task: %w", err)
	}

	return dp.showTask(task)
}

func (dp *DisplayPanel) ShowTaskByID(taskID int64) (database.Task, error) {
	task, err := database.GetTask(int64(taskID))
	if err != nil {
		return database.Task{}, fmt.Errorf("could not retrieve task [%d]: %w", taskID, err)
	}

	return dp.showTask(task)
}

type scrollDirection int

const (
	ScrollDown scrollDirection = iota
	ScrollUp
)

func (dp *DisplayPanel) ScrollTable(direction scrollDirection) error {
	table, ok := dp.frame.GetPrimitive().(*tview.Table)
	if !ok {
		return errors.New("task list is not available")
	}
	offset, _ := table.GetOffset()

	if direction == ScrollDown {
		// determine whether to fetch the next page
		rowCount := table.GetRowCount()
		_, _, _, height := table.GetInnerRect()
		if offset+height >= rowCount {
			return dp.GetNextPage()
		}
		table.SetOffset(offset+1, 0)
	} else {
		// determine whether to fetch the previous page
		if offset == 0 {
			return dp.getPreviousPage()
		}
		table.SetOffset(offset-1, 0)
	}
	return nil
}

func (dp *DisplayPanel) ScrollTask(direction scrollDirection) error {
	textView, ok := dp.frame.GetPrimitive().(*tview.TextView)
	if !ok {
		return errors.New("task is not available")
	}
	offset, _ := textView.GetScrollOffset()

	if direction == ScrollDown {
		offset, _ := textView.GetScrollOffset()
		textView.ScrollTo(offset+1, 0)
	} else if offset != 0 {
		textView.ScrollTo(offset-1, 0)
	}
	return nil
}
