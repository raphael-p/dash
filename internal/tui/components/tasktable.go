package components

import (
	"fmt"
	"time"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/raphael-p/dash/pkg/tcellutil"
	"github.com/raphael-p/dash/pkg/tviewcomponents/tableutil"
	"github.com/rivo/tview"
)

func NewTaskTable(tasks []database.Task) *tview.Table {
	table := tview.NewTable()

	for row, task := range tasks {
		var ts time.Time
		if task.CompletedAt.Valid {
			ts = task.CompletedAt.Time
		} else if task.PausedAt.Valid {
			ts = task.PausedAt.Time
		} else {
			ts = task.CreatedAt
		}

		idCell := tview.NewTableCell(fmt.Sprintf("[%d]", task.ID))
		nameCell := tview.NewTableCell(task.Name).SetExpansion(1).SetReference(task.Name)
		timeAgoCell := tview.NewTableCell(strutil.RelativeTime(ts)).SetTextColor(tcellutil.ColourCodeByAge(ts))
		table.SetCell(row, 0, idCell).SetCell(row, 1, nameCell).SetCell(row, 2, timeAgoCell)
	}

	tableutil.EnableDynamicTruncation(table, 1, 0, 2)

	return table
}
