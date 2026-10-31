package panels

import (
	"fmt"

	"github.com/rivo/tview"
)

type InfoPanel struct {
	textView *tview.TextView
}

func NewInfoPanel() *InfoPanel {
	textView := tview.NewTextView().SetScrollable(true).SetDynamicColors(true)
	textView.SetBorderPadding(1, 1, 2, 2)
	return &InfoPanel{textView}
}

func (ip *InfoPanel) GetPrimitive() tview.Primitive {
	return ip.textView
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
