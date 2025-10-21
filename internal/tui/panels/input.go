package panels

import (
	"fmt"
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputPanel struct {
	panel         *tview.Frame
	timerStart    time.Time
	timerDuration time.Duration
}

func NewInputPanel() *InputPanel {
	return &InputPanel{tview.NewFrame(tview.NewTextArea()), time.Time{}, time.Minute * 20}
}

func (ip *InputPanel) GetPanel() *tview.Frame {
	return ip.panel
}

func (ip *InputPanel) Set(title string, primitive tview.Primitive) {
	ip.panel.SetTitle(" " + title + " ")
	ip.panel.SetBorder(true).SetBorderColor(tcell.ColorLimeGreen)
	ip.panel.SetPrimitive(primitive)
}

func (ip *InputPanel) RestartTimer() {
	ip.timerStart = time.Now()
}

func (ip *InputPanel) GetMinutesRemaining() string {
	timeElapsed := time.Since(ip.timerStart)
	timeRemaining := ip.timerDuration - timeElapsed
	if timeRemaining < 0 || timeRemaining > ip.timerDuration {
		return "0"
	} else {
		return fmt.Sprintf("%d", int(math.Ceil(timeRemaining.Minutes())))
	}
}
