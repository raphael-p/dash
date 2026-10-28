package components

import "github.com/rivo/tview"

const headerFooterPadding int = 1

type Frame struct {
	*tview.Frame
}

func NewFrame() *Frame {
	frame := tview.NewFrame(nil).SetBorders(1, 1, headerFooterPadding, headerFooterPadding, 2, 2)
	frame.SetBorder(true)
	return &Frame{frame}
}

func (f *Frame) GetAvailableRows() int {
	_, _, _, height := f.GetInnerRect()
	headerAndFooterRows := (1 + 2*headerFooterPadding) * 2
	return height - headerAndFooterRows
}
