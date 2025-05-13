package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputPanel struct {
	panel *tview.Frame
}

func NewInputPanel() *InputPanel {
	return &InputPanel{tview.NewFrame(tview.NewTextArea())}
}

func (ip *InputPanel) GetPanel() *tview.Frame {
	return ip.panel
}

func (ip *InputPanel) Set(title string, primitive tview.Primitive) {
	ip.panel.SetTitle(" " + title + " ")
	ip.panel.SetBorder(true).SetBorderColor(tcell.ColorLimeGreen)
	ip.panel.SetPrimitive(primitive)
}
