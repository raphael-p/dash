package panels

import (
	"fmt"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/rivo/tview"
)

type taskPage struct {
	lastId  int64
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

func (dp *DisplayPanel) GetCurrentPage() error {
	if dp.currentPageIndex == 0 {
		return dp.ListTasks(0, 0)
	}
	return dp.ListTasks(
		dp.pages[dp.currentPageIndex-1].lastId,
		dp.currentPageIndex,
	)
}

func (dp *DisplayPanel) GetPreviousPage() error {
	if dp.currentPageIndex == 0 {
		return nil // intentional noop
	}

	pageIndex := dp.currentPageIndex - 1

	// get last ID of previous page
	var fromId int64 = 0
	if pageIndex > 0 {
		fromId = dp.pages[pageIndex-1].lastId
	}

	err := dp.ListTasks(fromId, pageIndex)
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
	err := dp.ListTasks(currentPage.lastId, pageIndex)
	if err == nil {
		dp.currentPageIndex = pageIndex
		dp.GetPanel().ScrollToBeginning()
	}
	return err
}

func (dp *DisplayPanel) ListTasks(fromId int64, pageIndex uint) error {
	dp.panel.Clear()
	dp.panel.SetTitle(" Tasks ")

	tasks, hasNext, err := database.GetTasksPaginated(fromId)
	if err != nil {
		return fmt.Errorf("could not retrieve tasks: %s", err)
	}

	if len(tasks) == 0 {
		fmt.Fprintln(dp.panel, "No tasks. Good job :)")
		return nil
	}

	// store pagination metadata
	dp.pages[pageIndex] = taskPage{tasks[len(tasks)-1].Id, hasNext}

	for idx, task := range tasks {
		fmt.Fprintf(dp.panel, "[%d] %s", task.Id, task.Name)
		if idx+1 < len(tasks) {
			fmt.Fprint(dp.panel, "\n")
		}
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
		return dp.GetPreviousPage()
	}

	dp.GetPanel().ScrollTo(offset-1, 0)
	return nil
}
