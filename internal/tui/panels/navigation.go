package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/tui/components"
	"github.com/rivo/tview"
)

type NavigationPanel struct {
	frame    *components.Frame
	textView *tview.TextView
}

func NewNavigationPanel(borderColour, headerFooterColour tcell.Color) *NavigationPanel {
	textView := tview.NewTextView().SetDynamicColors(true)
	frame := components.NewFrame(headerFooterColour)
	frame.SetBorderColor(borderColour)
	return &NavigationPanel{frame, textView}
}

func (nav *NavigationPanel) GetPrimitive() tview.Primitive {
	return nav.frame
}

func (nav *NavigationPanel) SetPrimitive(primitime tview.Primitive) *NavigationPanel {
	nav.frame.SetPrimitive(primitime)
	return nav
}

func (nav *NavigationPanel) SetText(text string) *NavigationPanel {
	nav.frame.SetPrimitive(nav.textView)
	nav.textView.SetText(text)
	return nav
}

func (nav *NavigationPanel) SetTitle(title string) *NavigationPanel {
	nav.frame.SetTitle(" " + title + " ")
	return nav
}
