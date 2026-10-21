package commandline

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CommandLine struct {
	flex     *tview.Flex
	textArea *tview.TextArea
}

func New(chevronColour tcell.Color) *CommandLine {
	chevron := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf("[%s::b]>[-:-:-] ", chevronColour))
	textArea := tview.NewTextArea()
	textArea.SetPlaceholderStyle(textArea.GetPlaceholderStyle().Foreground(tcell.ColorGray))
	flex := tview.NewFlex().AddItem(tview.NewBox(), 2, 0, false).AddItem(chevron, 2, 0, false).AddItem(textArea, 0, 1, true)
	return &CommandLine{flex, textArea}
}

func (cl *CommandLine) SetChangedFunc(handler func()) *CommandLine {
	cl.textArea.SetChangedFunc(handler)
	return cl
}

func (cl *CommandLine) SetCharLimit(charLimit int) *CommandLine {
	cl.textArea.SetMaxLength(charLimit)
	return cl
}

func (cl *CommandLine) SetPlaceholder(placeholder string) *CommandLine {
	cl.textArea.SetPlaceholder(placeholder)
	return cl
}

func (cl *CommandLine) GetInput() string {
	return cl.textArea.GetText()
}

func (cl *CommandLine) HasInput() bool {
	return cl.textArea.GetText() != ""
}

func (cl *CommandLine) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box {
	return cl.flex.SetInputCapture(capture)
}

func (cl *CommandLine) GetPrimitive() tview.Primitive {
	return cl.flex
}

func (cl *CommandLine) Reset() *CommandLine {
	cl.SetChangedFunc(nil).SetPlaceholder("").SetCharLimit(0)
	cl.textArea.SetText("", true)
	return cl
}
