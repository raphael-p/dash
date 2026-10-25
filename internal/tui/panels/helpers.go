package panels

import "github.com/rivo/tview"

func newFrame() *tview.Frame {
	frame := tview.NewFrame(nil).SetBorders(1, 1, 1, 1, 2, 2)
	frame.SetBorder(true)
	return frame
}
