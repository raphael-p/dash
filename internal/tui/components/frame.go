package components

import (
	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/rivo/tview"
)

const headerFooterPadding int = 1

type Frame struct {
	*tview.Frame
}

func NewFrame() *Frame {
	frame := tview.NewFrame(nil).SetBorders(1, 1, headerFooterPadding, headerFooterPadding, 2, 2)
	frame.SetBorder(true)
	return &Frame{frame}
}

func (f *Frame) AddHeader(header string) {
	f.AddText(header, true, tview.AlignLeft, config.C.SecondaryColour)
}

func (f *Frame) AddFooter(footer string) {
	f.AddText(footer, false, tview.AlignLeft, config.C.SecondaryColour)
}

func (f *Frame) AddDateFooter(date, label string) {
	if date != "" {
		f.AddFooter(strutil.RightPad(label, 14) + date)
	}
}

func (f *Frame) GetAvailableRows() int {
	_, _, _, height := f.GetInnerRect()
	headerAndFooterRows := (1 + 2*headerFooterPadding) * 2
	return height - headerAndFooterRows
}
