package commandline

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CommandLine struct {
	flex          *tview.Flex
	textArea      *tview.TextArea
	chevronColour tcell.Color
	input         string
	onInputChange func(string)
	placeholder   string
}

func New(chevronColour tcell.Color, placeholder string) *CommandLine {
	chevron := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf("[%s::b]>[-:-:-] a", chevronColour))
	textArea := tview.NewTextArea().SetPlaceholder(" " + placeholder)
	textArea.SetPlaceholderStyle(textArea.GetPlaceholderStyle().Foreground(tcell.ColorGray))
	flex := tview.NewFlex().AddItem(chevron, 1, 0, false).AddItem(textArea, 0, 1, false)

	return &CommandLine{flex, textArea, 0, "", func(string) {}, ""}
}

func (cl *CommandLine) GetInput() string {
	return cl.input
}

func (cl *CommandLine) SetInput(input string) {
	cl.input = input
	if input == "" {
		cl.textArea.SetText("", true)
	} else {
		cl.textArea.SetText(" "+cl.input, true)
	}
	cl.onInputChange(cl.input)
}

func (cl *CommandLine) SetOnInputChange(onInputChange func(newInput string)) {
	cl.onInputChange = onInputChange
}

func (cl *CommandLine) GetPrimitive() tview.Primitive {
	return cl.flex
}
