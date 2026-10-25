package panels

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newFrame() *tview.Frame {
	frame := tview.NewFrame(nil).SetBorders(1, 1, 1, 1, 2, 2)
	frame.SetBorder(true)
	return frame
}

func colorCodeByAge(t time.Time) tcell.Color {
	d := time.Since(t)
	day := 24 * time.Hour
	switch {
	case d < day:
		return tcell.Color220
	case d < day*7:
		return tcell.Color214
	case d < day*15:
		return tcell.Color208
	case d < day*30:
		return tcell.Color202
	default:
		return tcell.Color160
	}
}

func formatTimeAgo(t time.Time) string {
	d := time.Since(t)
	day := 24 * time.Hour
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < day:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d < 7*day:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	case d < 365*day:
		return fmt.Sprintf("%dw", int(d.Hours()/(24*7)))
	default:
		return fmt.Sprintf("%dy", int(d.Hours()/(24*365)))
	}
}
