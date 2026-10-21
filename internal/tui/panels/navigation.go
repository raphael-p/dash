package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NavigationPanel struct {
	frame    *tview.Frame
	textView *tview.TextView
}

func NewNavigationPanel() *NavigationPanel {
	textView := tview.NewTextView().SetDynamicColors(true)
	frame := tview.NewFrame(nil)
	frame.SetBorder(true).SetBorderColor(tcell.ColorLimeGreen).SetBorderPadding(1, 1, 2, 2)
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
