package panels

import (
	"fmt"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/stringpad"
	"github.com/rivo/tview"
)

type taskPage struct {
	lastID  int64
	hasNext bool
}

type DisplayPanel struct {
	panel            *tview.TextView
	pages            map[uint]taskPage
	currentPageIndex uint
}

func NewDisplayPanel() *DisplayPanel {
	panel := tview.NewTextView().SetWordWrap(true)
	panel.SetBorder(true).SetBorderPadding(1, 1, 2, 2)
	return &DisplayPanel{panel, map[uint]taskPage{}, 0}
}

func (dp *DisplayPanel) GetPanel() *tview.TextView {
	return dp.panel
}

func (dp *DisplayPanel) ResetPagination() {
	dp.currentPageIndex = 0
}

func (dp *DisplayPanel) GetCurrentPage() error {
	if dp.currentPageIndex == 0 {
		return dp.listTasks(0, 0)
	}
	return dp.listTasks(
		dp.pages[dp.currentPageIndex-1].lastID,
		dp.currentPageIndex,
	)
}

func (dp *DisplayPanel) getPreviousPage() error {
	if dp.currentPageIndex == 0 {
		return nil // intentional noop
	}

	pageIndex := dp.currentPageIndex - 1

	// get last ID of previous page
	var fromID int64 = 0
	if pageIndex > 0 {
		fromID = dp.pages[pageIndex-1].lastID
	}

	err := dp.listTasks(fromID, pageIndex)
	if err == nil {
		dp.currentPageIndex = pageIndex
		dp.GetPanel().ScrollToEnd()
	}
	return err
}

func (dp *DisplayPanel) GetNextPage() error {
	currentPage := dp.pages[dp.currentPageIndex]
	if !currentPage.hasNext {
		return nil // intentional noop
	}

	pageIndex := dp.currentPageIndex + 1
	err := dp.listTasks(currentPage.lastID, pageIndex)
	if err == nil {
		dp.currentPageIndex = pageIndex
		dp.GetPanel().ScrollToBeginning()
	}
	return err
}

func (dp *DisplayPanel) listTasks(fromID int64, pageIndex uint) error {
	dp.panel.Clear()
	dp.panel.SetTitle(fmt.Sprintf(" Tasks (page %d) ", pageIndex+1))

	tasks, hasNext, err := database.GetTasksPaginated(fromID)
	if err != nil {
		return fmt.Errorf("could not retrieve tasks: %s", err)
	}

	if len(tasks) == 0 {
		fmt.Fprintln(dp.panel, "No tasks. Good job :)")
		return nil
	}

	// store pagination metadata
	dp.pages[pageIndex] = taskPage{tasks[len(tasks)-1].Id, hasNext}

	if pageIndex > 0 {
		fmt.Fprintf(dp.panel, "\n↑ show previous page\n\n")
	}
	for idx, task := range tasks {
		fmt.Fprint(
			dp.panel,
			stringpad.RightPad(fmt.Sprintf("[%d]", task.Id), 8),
			task.Name,
		)
		if idx+1 < len(tasks) {
			fmt.Fprint(dp.panel, "\n")
		}
	}
	if hasNext {
		fmt.Fprintf(dp.panel, "\n\n↓ show next page\n")
	}
	return nil
}

func (dp *DisplayPanel) showTask(task database.Task) (database.Task, error) {
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
	return task, nil
}

func (dp *DisplayPanel) ShowTopTask() (database.Task, error) {
	task, err := database.GetTopTask()
	if err != nil {
		return database.Task{}, fmt.Errorf("could not retrieve top task: %s", err)
	}

	return dp.showTask(task)
}

func (dp *DisplayPanel) ShowTaskById(taskID int64) (database.Task, error) {
	task, err := database.GetTask(int64(taskID))
	if err != nil {
		return database.Task{}, fmt.Errorf("could not retrieve task %d: %s", taskID, err)
	}

	return dp.showTask(task)
}

func (dp *DisplayPanel) ScrollDown() error {
	offset, _ := dp.GetPanel().GetScrollOffset()

	// determine whether to fetch the next page
	lineCount := dp.GetPanel().GetWrappedLineCount()
	_, _, _, height := dp.GetPanel().GetInnerRect()
	if offset+height >= lineCount {
		return dp.GetNextPage()
	}

	dp.GetPanel().ScrollTo(offset+1, 0)
	return nil
}

func (dp *DisplayPanel) ScrollUp() error {
	offset, _ := dp.GetPanel().GetScrollOffset()

	// determine whether to fetch the previous page
	if offset == 0 {
		return dp.getPreviousPage()
	}

	dp.GetPanel().ScrollTo(offset-1, 0)
	return nil
}
