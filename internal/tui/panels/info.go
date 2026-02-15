package panels

import (
	"fmt"

	"github.com/rivo/tview"
)

type InfoPanel struct {
	panel         *tview.Frame
	textView      *tview.TextView
	input         string
	onInputChange func(string)
}

func NewInfoPanel() *InfoPanel {
	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
	panel := tview.NewFrame(textView)
	return &InfoPanel{panel, textView, "", func(string) {}}
}

func (ip *InfoPanel) SetOnInputChange(onInputChange func(newInput string)) {
	ip.onInputChange = onInputChange
}

func (ip *InfoPanel) GetPanel() *tview.Frame {
	return ip.panel
}

func (ip *InfoPanel) Warn(err string) {
	ip.clearInput()
	ip.textView.SetText(fmt.Sprint("[yellow]⚠️ ", err))
}

func (ip *InfoPanel) Error(err error) {
	ip.clearInput()
	ip.textView.SetText(fmt.Sprint("[red]❌ ", err))
}

func (ip *InfoPanel) Info(info string) {
	ip.clearInput()
	ip.textView.SetText(fmt.Sprint("❕ ", info))
}

func (ip *InfoPanel) SetInput(input string) {
	ip.input = input
	if input != "" {
		ip.textView.SetText(fmt.Sprint("> ", ip.input))
	} else {
		ip.textView.Clear()
	}
	ip.onInputChange(ip.input)
}

func (ip *InfoPanel) GetInput() string {
	return ip.input
}

func (ip *InfoPanel) clearInput() {
	if ip.input != "" {
		ip.input = ""
		ip.onInputChange("")
	}
}

func (ip *InfoPanel) Clear() {
	ip.input = ""
	ip.textView.Clear()
}
