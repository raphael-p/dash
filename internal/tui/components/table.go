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

type Table struct {
	*tview.Table
}

func NewTable(tasks []database.Task) *Table {
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
		nameEscaped := tview.Escape(task.Name)
		nameCell := tview.NewTableCell(nameEscaped).SetExpansion(1).SetReference(nameEscaped)
		timeAgoCell := tview.NewTableCell(strutil.RelativeTime(ts)).SetTextColor(tcellutil.ColourCodeByAge(ts))
		table.SetCell(row, 0, idCell).SetCell(row, 1, nameCell).SetCell(row, 2, timeAgoCell)
	}

	tableutil.EnableDynamicTruncation(table, 1, 0, 2)

	return &Table{table}
}

func (t *Table) Scroll(direction ScrollDirection) bool {
	offset, _ := t.GetOffset()

	if direction == ScrollDown {
		// determine whether to fetch the next page
		rowCount := t.GetRowCount()
		_, _, _, height := t.GetInnerRect()
		if offset+height >= rowCount {
			return true
		}
		t.SetOffset(offset+1, 0)
	} else {
		// determine whether to fetch the previous page
		if offset == 0 {
			return true
		}
		t.SetOffset(offset-1, 0)
	}
	return false
}
