package panels

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/stringpad"
	"github.com/rivo/tview"
)

func newFrame() *tview.Frame {
	frame := tview.NewFrame(nil).SetBorders(1, 1, 1, 1, 2, 2)
	frame.SetBorder(true)
	return frame
}

func makeTaskTable(tasks []database.Task, isCompleted bool) *tview.Table {
	table := tview.NewTable()

	for row, task := range tasks {
		var ts time.Time
		if isCompleted && task.CompletedAt.Valid {
			ts = task.CompletedAt.Time
		} else {
			ts = task.CreatedAt
		}

		idCell := tview.NewTableCell(fmt.Sprintf("[%d]", task.ID))
		nameCell := tview.NewTableCell(task.Name).SetExpansion(1).SetReference(task.Name)
		timeAgoCell := tview.NewTableCell(formatTimeAgo(ts)).SetTextColor(colorCodeByAge(ts))
		table.SetCell(row, 0, idCell).SetCell(row, 1, nameCell).SetCell(row, 2, timeAgoCell)
	}

	enableDynamicTruncation(table, 1, 0, 2)

	return table
}

// Truncate targetCol text on table size changes so fixedCols are always displayed
func enableDynamicTruncation(table *tview.Table, targetCol int, fixedCols ...int) {
	var lastWidth int

	table.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		_, _, tabeWidth, _ := table.GetInnerRect()
		if tabeWidth == lastWidth || tabeWidth <= 0 {
			return x, y, width, height
		}
		lastWidth = tabeWidth

		fixedWidth := 0
		for _, col := range fixedCols {
			_, _, colWidth := table.GetCell(0, col).GetLastPosition()
			fixedWidth += colWidth
		}
		separatorWidth := table.GetColumnCount()
		remainingWidth := max(tabeWidth-fixedWidth-separatorWidth, 4)

		for r := 0; r < table.GetRowCount(); r++ {
			cell := table.GetCell(r, targetCol)
			if cell == nil {
				continue
			}
			text, ok := cell.GetReference().(string)
			if !ok || text == "" {
				text = cell.Text
			}
			cell.SetText(stringpad.TruncateString(text, remainingWidth))
		}

		return x, y, width, height
	})
}

func colorCodeByAge(t time.Time) tcell.Color {
	d := time.Since(t)
	day := 24 * time.Hour
	switch {
	case d < day:
		return tcell.Color220
	case d < day*7:
		return tcell.Color214
	case d < day*15:
		return tcell.Color208
	case d < day*30:
		return tcell.Color202
	default:
		return tcell.Color160
	}
}

func formatTimeAgo(t time.Time) string {
	d := time.Since(t)
	day := 24 * time.Hour
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < day:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d < 7*day:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	case d < 365*day:
		return fmt.Sprintf("%dw", int(d.Hours()/(24*7)))
	default:
		return fmt.Sprintf("%dy", int(d.Hours()/(24*365)))
	}
}
