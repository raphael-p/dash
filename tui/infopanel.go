package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type infoPanel struct {
	panel    *tview.Frame
	textView *tview.TextView
}

func makeInfoPanel() *infoPanel {
	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
	panel := tview.NewFrame(textView)
	return &infoPanel{panel, textView}
}

func (ip *infoPanel) warning(err any) {
	ip.textView.SetText(fmt.Sprint("[yellow]⚠️ ", err))
}

func (ip *infoPanel) error(err any) {
	ip.textView.SetText(fmt.Sprint("[red]❌ ", err))
}

func (ip *infoPanel) message(info string) {
	ip.textView.SetText(fmt.Sprint("❕ ", info))
}

func (ip *infoPanel) clear() {
	ip.textView.Clear()
}
