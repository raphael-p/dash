package commandline

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/config"
	"github.com/rivo/tview"
)

type CommandLine struct {
	flex     *tview.Flex
	textArea *tview.TextArea
}

func New() *CommandLine {
	chevron := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf("[%s::b]>[-:-:-] ", config.C.HighlightColour))
	textArea := tview.NewTextArea()
	textArea.SetPlaceholderStyle(textArea.GetPlaceholderStyle().Foreground(config.C.SecondaryColour))
	flex := tview.NewFlex().AddItem(tview.NewBox(), 2, 0, false).AddItem(chevron, 2, 0, false).AddItem(textArea, 0, 1, true)
	return &CommandLine{flex, textArea}
}

func (cl *CommandLine) SetChangedFunc(handler func()) *CommandLine {
	cl.textArea.SetChangedFunc(handler)
	return cl
}

func (cl *CommandLine) Clear() *CommandLine {
	cl.textArea.SetText("", true)
	cl.SetChangedFunc(nil)
	return cl
}

func (cl *CommandLine) GetInput() string {
	return cl.textArea.GetText()
}

func (cl *CommandLine) HasInput() bool {
	return cl.textArea.GetText() != ""
}

func (cl *CommandLine) SetInput(input string) *CommandLine {
	cl.textArea.SetText(input, true)
	return cl
}

func (cl *CommandLine) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box {
	return cl.flex.SetInputCapture(capture)
}

func (cl *CommandLine) GetPrimitive() tview.Primitive {
	return cl.flex
}

func (cl *CommandLine) Set(charLimit int, placeholder string, capture func(event *tcell.EventKey) *tcell.EventKey) *CommandLine {
	cl.flex.SetInputCapture(capture)
	cl.textArea.SetMaxLength(charLimit).SetPlaceholder(placeholder)
	return cl
}
