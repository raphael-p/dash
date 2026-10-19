package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputPanel struct {
	frame    *tview.Frame
	getInput func() string
	setInput func(string)
}

func NewInputPanel() *InputPanel {
	return &InputPanel{tview.NewFrame(tview.NewTextArea()), func() string { return "" }, func(string) {}}
}

func (ip *InputPanel) GetInput() string {
	return ip.getInput()
}

func (ip *InputPanel) SetInput(input string) {
	ip.setInput(input)
}

func (ip *InputPanel) GetPrimitive() *tview.Frame {
	return ip.frame
}

func (ip *InputPanel) Set(title string, primitive tview.Primitive, getInput func() string, setInput func(string)) {
	ip.frame.SetTitle(" " + title + " ")
	ip.frame.SetBorder(true).SetBorderColor(tcell.ColorLimeGreen)
	ip.frame.SetPrimitive(primitive)
	ip.getInput = getInput
	ip.setInput = setInput
}
