package tviewcomponents

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/rivo/tview"
)

// EnableDynamicTruncation fits a plain-text target column around fixed columns.
// Store the raw target text in each cell's string reference before enabling it.
// This helper is for borderless tables with a single-cell column separator.
func EnableDynamicTruncation(table *tview.Table, targetCol int, fixedCols ...int) {
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
			cell.SetText(strutil.UnicodeTruncate(text, remainingWidth))
		}

		return x, y, width, height
	})
}
