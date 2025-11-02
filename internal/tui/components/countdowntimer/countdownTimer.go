package countdowntimer

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const DURATION = time.Minute * 20
const REFRESH_INTERVAL = time.Second * 2

type countdownTimer struct {
	Layout      *tview.Flex
	countdown   *tview.TextView
	description *tview.TextView
	startTime   atomic.Int64
}

var instance countdownTimer

var Instance = sync.OnceValue(func() *countdownTimer {
	if instance.Layout != nil {
		return &instance
	}

	countdown := tview.NewTextView()
	countdown.SetTextColor(tcell.ColorLimeGreen)
	countdown.SetBorderPadding(1, 0, 2, 2)
	instance.countdown = countdown

	description := tview.NewTextView().SetDynamicColors(true)
	description.SetBorderPadding(0, 1, 2, 2)
	instance.description = description

	instance.Layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(countdown, 3, 0, false).
		AddItem(description, 0, 1, false)

	return &instance
})

func (t *countdownTimer) SetDescription(text string) {
	t.description.SetText(text)
}

func (t *countdownTimer) getMinutesRemaining() string {
	startTime := time.Unix(t.startTime.Load(), 0)
	timeElapsed := time.Since(startTime)
	timeRemaining := DURATION - timeElapsed
	if timeRemaining < 0 || timeRemaining > DURATION {
		return "0"
	} else {
		return fmt.Sprintf("%d", int(math.Ceil(timeRemaining.Minutes())))
	}
}

func (t *countdownTimer) Reset(app *tview.Application) {
	t.startTime.Store(0)
	t.Start(app)
}

func (t *countdownTimer) Start(app *tview.Application) {
	if !t.startTime.CompareAndSwap(0, time.Now().Unix()) {
		return // avoid dangling threads
	}

	go (func(app *tview.Application, t *countdownTimer) {
		for {
			minutesRemaining := t.getMinutesRemaining()
			if minutesRemaining == "0" {
				break
			}

			app.QueueUpdateDraw(func() {
				t.countdown.SetText(fmt.Sprintf("%s minutes remaining", minutesRemaining))
			})
			time.Sleep(REFRESH_INTERVAL)
		}

		app.QueueUpdateDraw(func() {
			t.countdown.SetText("dash complete. restart when ready.")
		})

		t.startTime.Store(0)
	})(app, t)
}
