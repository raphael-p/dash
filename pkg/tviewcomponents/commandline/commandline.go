package commandline

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CommandLine struct {
	chevronColour tcell.Color
	input         string
	onInputChange func(string)
	placeholder   string
	textView      *tview.TextView
}

func New() *CommandLine {
	textView := tview.NewTextView().SetDynamicColors(true)
	return &CommandLine{0, "", func(string) {}, "\n", textView}
}

func (cl *CommandLine) SetChevronColour(colour tcell.Color) *CommandLine {
	cl.chevronColour = colour
	return cl
}

func (cl *CommandLine) SetPlaceholder(placeholder string) *CommandLine {
	cl.placeholder = placeholder
	return cl
}

func (cl *CommandLine) GetInput() string {
	return cl.input
}

func (cl *CommandLine) SetInput(input string) {
	cl.input = input
	if input != "" {
		cl.textView.SetText(fmt.Sprintf("[%s::b]>[-:-:-] %s", cl.chevronColour, cl.input))
	} else {
		cl.textView.SetText(fmt.Sprintf("[%s::b]>[-:-:-] [%s::b]%s[-:-:-]", cl.chevronColour, tcell.ColorGray, cl.placeholder))
	}
	cl.onInputChange(cl.input)
}

func (cl *CommandLine) SetOnInputChange(onInputChange func(newInput string)) {
	cl.onInputChange = onInputChange
}

func (cl *CommandLine) GetTextView() *tview.TextView {
	return cl.textView
}
