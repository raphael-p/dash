package components

import (
	"github.com/rivo/tview"
)

// heuristic to avoid truncation
const (
	minArea, minHeight, minWidth int = 1200, 16, 40
	reducedHeight, reducedWidth  int = 27, 80
	infoHeight, spacerHeight     int = 5, 1
)

type Grid struct {
	*tview.Grid
	display, info, nav, command tview.Primitive
	spacer                      tview.Primitive
	width, height               int
}

func NewGrid(display, info, nav, command tview.Primitive) *Grid {
	g := &Grid{
		Grid:    tview.NewGrid(),
		display: display, info: info, nav: nav, command: command,
		spacer: tview.NewBox(),
	}
	g.layout(reducedWidth, reducedHeight)
	return g
}

func (g *Grid) SetRect(x, y, width, height int) {
	g.Grid.SetRect(x, y, width, height)
	if width != g.width || height != g.height {
		g.layout(width, height)
	}
}

func (g *Grid) layout(width, height int) {
	g.width, g.height = width, height
	g.Clear().SetOffset(0, 0)
	switch {
	case width < minWidth || height < minHeight || width*height < minArea:
		g.SetRows(-1, 1).SetColumns(-1).
			AddItem(g.display, 0, 0, 1, 1, 0, 0, false).
			AddItem(g.command, 1, 0, 1, 1, 0, 0, true)
	case height < reducedHeight, width < reducedWidth:
		g.SetRows(-3, -1).SetColumns(-2, -1).
			AddItem(g.display, 0, 0, 1, 1, 0, 0, false).
			AddItem(g.nav, 0, 1, 1, 1, 0, 0, false).
			AddItem(g.command, 1, 0, 1, 2, 0, 0, true)
	default:
		g.SetRows(-2, infoHeight, spacerHeight, -1).SetColumns(-3, -1).
			AddItem(g.display, 0, 0, 2, 1, 0, 0, false).
			AddItem(g.nav, 0, 1, 1, 1, 0, 0, false).
			AddItem(g.info, 1, 1, 1, 1, 0, 0, false).
			AddItem(g.spacer, 2, 0, 1, 2, 0, 0, false).
			AddItem(g.command, 3, 0, 1, 2, 0, 0, true)
	}
}
