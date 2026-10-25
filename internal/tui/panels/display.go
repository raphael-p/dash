package panels

import (
	"errors"
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
	frame        *tview.Frame
	activeTab    *displayPanelTab
	completedTab *displayPanelTab
	mode         database.TaskMode
}

func NewDisplayPanel() *DisplayPanel {
	frame := tview.NewFrame(nil).SetBorders(2, 2, 0, 0, 4, 4)
	frame.SetBorder(true)
	return &DisplayPanel{frame, newDisplayPanelTab(), newDisplayPanelTab(), database.ActiveTasks}
}

func (dp *DisplayPanel) GetPrimitive() tview.Primitive {
	return dp.frame
}

func (dp *DisplayPanel) getTable() (*tview.Table, error) {
	if table, ok := dp.frame.GetPrimitive().(*tview.Table); ok {
		return table, nil
	}
	return nil, errors.New("task list is not available")
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
	return dp.listTasks(dp.getTab().currentPageIndex)
}

// TODO: fix scrolling to end when going back a page
func (dp *DisplayPanel) getPreviousPage() error {
	table, err := dp.getTable()
	if err != nil {
		return err
	}
	tab := dp.getTab()
	if tab.currentPageIndex == 0 {
		return nil // intentional noop
	}

	pageIndex := tab.currentPageIndex - 1
	err = dp.listTasks(pageIndex)
	if err == nil {
		tab.currentPageIndex = pageIndex
		table.ScrollToEnd()
	}
	return err
}

func (dp *DisplayPanel) GetNextPage() error {
	table, err := dp.getTable()
	if err != nil {
		return err
	}
	tab := dp.getTab()
	currentPage := tab.pages[tab.currentPageIndex]
	if !currentPage.hasNext {
		return nil // intentional noop
	}

	pageIndex := tab.currentPageIndex + 1
	err = dp.listTasks(pageIndex)
	if err == nil {
		tab.currentPageIndex = pageIndex
		table.ScrollToBeginning()
	}
	return err
}

func (dp *DisplayPanel) SetTitle(title string) *DisplayPanel {
	dp.frame.SetTitle(title)
	return dp
}

func (dp *DisplayPanel) listTasks(pageIndex uint) error {
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
	}

	tasks, hasNext, err := database.GetTasksPaginated(dp.mode, cursor)
	if err != nil {
		return fmt.Errorf("could not retrieve tasks: %w", err)
	}

	if len(tasks) == 0 {
		pages[pageIndex] = taskPage{}
		dp.frame.SetPrimitive(tview.NewTextView().SetText("No tasks. Good job :)"))
		return nil
	}

	// store pagination metadata
	pages[pageIndex] = taskPage{tasks[len(tasks)-1], hasNext}

	table := tview.NewTable()
	defer dp.frame.SetPrimitive(table)
	if pageIndex > 0 {
		table.SetCellSimple(1, 0, "").SetCellSimple(1, 1, "↑ show previous page").SetCellSimple(2, 0, "")
	}
	rowCount := table.GetRowCount()
	for taskIndex, task := range tasks {
		i := taskIndex + rowCount
		table.SetCellSimple(i, 0, fmt.Sprintf("[%d]", task.ID)).SetCellSimple(i, 1, task.Name)
		// TODO: format properly !!!
		// if isCompleted {
		// 	table.SetCellSimple(i, 2, task.CompletedAt.Time.Format("Jan 2"))
		// }
	}

	rowCount = table.GetRowCount()
	if hasNext {
		table.SetCellSimple(rowCount+1, 0, "").SetCellSimple(rowCount+1, 1, "↓ show next page").SetCellSimple(rowCount+2, 0, "")
	}
	return nil
}

func (dp *DisplayPanel) showTask(task database.Task) (database.Task, error) {
	dp.SetTitle(fmt.Sprintf("Task %d: %s", task.ID, task.Name))
	textView := tview.NewTextView()
	defer dp.frame.SetPrimitive(textView)
	if task.Description != "" {
		fmt.Fprintf(textView, "%s\n\n", task.Description)
	}

	createdAt, updatedAt, completedAt := task.GetFormattedTimes()

	labelise := func(text string) string { return stringpad.RightPad(text, 20) }
	if completedAt != "" {
		fmt.Fprintln(textView, labelise("completed at"), completedAt)
	}
	fmt.Fprintln(textView, labelise("created at"), createdAt)
	fmt.Fprintln(textView, labelise("last updated at"), updatedAt)
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

// TODO: implement scroll for text view as well (necessary for big description)
func (dp *DisplayPanel) ScrollDown() error {
	table, err := dp.getTable()
	if err != nil {
		return err
	}
	offset, _ := table.GetOffset()

	// determine whether to fetch the next page
	rowCount := table.GetRowCount()
	_, _, _, height := table.GetInnerRect()
	if offset+height >= rowCount {
		return dp.GetNextPage()
	}

	table.SetOffset(offset+1, 0)
	return nil
}

func (dp *DisplayPanel) ScrollUp() error {
	table, err := dp.getTable()
	if err != nil {
		return err
	}
	offset, _ := table.GetOffset()

	// determine whether to fetch the previous page
	if offset == 0 {
		return dp.getPreviousPage()
	}

	table.SetOffset(offset-1, 0)
	return nil
}
