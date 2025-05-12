package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type inputPanel struct {
	panel *tview.Frame
}

func makeInputPanel() *inputPanel {
	return &inputPanel{tview.NewFrame(tview.NewTextArea())}
}

func (ip *inputPanel) set(title string, primitive tview.Primitive) {
	ip.panel.SetTitle(" " + title + " ")
	ip.panel.SetBorder(true).SetBorderColor(tcell.ColorLimeGreen)
	ip.panel.SetPrimitive(primitive)
}
