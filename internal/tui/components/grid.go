package components

import (
	"github.com/rivo/tview"
)

func NewGrid(
	display, info, nav, command tview.Primitive,
	minHeight, minWidth int,
) *tview.Grid {
	spacer := tview.NewBox()

	// minimal display
	grid := tview.NewGrid().SetRows(-2, 5, 1, -1).SetColumns(-3, -1).
		AddItem(display, 0, 0, 3, 2, 0, 0, false).
		AddItem(command, 3, 0, 1, 2, 0, 0, true)

	// short display
	grid.AddItem(display, 0, 0, 3, 1, 0, minWidth, false).
		AddItem(nav, 0, 1, 3, 1, 0, minWidth, false)

	// narrow display
	grid.AddItem(display, 0, 0, 2, 2, minHeight, 0, false)

	// regular display, added last so it wins ties with the short display
	grid.AddItem(display, 0, 0, 2, 1, minHeight, minWidth, false).
		AddItem(nav, 0, 1, 1, 1, minHeight, minWidth, false).
		AddItem(info, 1, 1, 1, 1, minHeight, minWidth, false).
		AddItem(spacer, 2, 0, 1, 2, minHeight, 0, false)

	return grid
}
