package panels

import (
	"fmt"

	"github.com/rivo/tview"
)

type InfoPanel struct {
	panel    *tview.Frame
	textView *tview.TextView
}

func NewInfoPanel() *InfoPanel {
	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
	panel := tview.NewFrame(textView)
	return &InfoPanel{panel, textView}
}

func (ip *InfoPanel) GetPanel() *tview.Frame {
	return ip.panel
}

func (ip *InfoPanel) Warn(err string) {
	ip.textView.SetText(fmt.Sprint("[yellow]⚠️ ", err))
}

func (ip *InfoPanel) Error(err error) {
	ip.textView.SetText(fmt.Sprint("[red]❌ ", err))
}

func (ip *InfoPanel) Info(info string) {
	ip.textView.SetText(fmt.Sprint("❕ ", info))
}

func (ip *InfoPanel) Clear() {
	ip.textView.Clear()
}
