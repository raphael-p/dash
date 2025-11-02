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

const COUNTDOWN_DURATION = time.Minute * 20
const REFRESH_INTERVAL = time.Second * 2
const START_MESSAGE_DURATION = time.Second * 2

type countdownTimer struct {
	Layout             *tview.Flex
	countdown          *tview.TextView
	description        *tview.TextView
	startMessage       string
	endMessage         string
	isRunning          atomic.Bool
	reinitialiseSignal chan struct{}
}

var instance countdownTimer
var once sync.Once

var Instance = func(startMessage, endMessage string) *countdownTimer {
	once.Do(func() {
		countdown := tview.NewTextView()
		countdown.SetTextColor(tcell.ColorLimeGreen)
		countdown.SetBorderPadding(1, 0, 2, 2)

		description := tview.NewTextView().SetDynamicColors(true)
		description.SetBorderPadding(0, 1, 2, 2)

		layout := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(countdown, 3, 0, false).
			AddItem(description, 0, 1, false)

		instance = countdownTimer{layout, countdown, description, startMessage, endMessage, atomic.Bool{}, make(chan struct{}, 1)}
	})
	return &instance
}

func (t *countdownTimer) SetDescription(text string) {
	t.description.SetText(text)
}

func getMinutesRemaining(startTime time.Time) string {
	timeElapsed := time.Since(startTime)
	timeRemaining := COUNTDOWN_DURATION - timeElapsed
	if timeRemaining < 0 || timeRemaining > COUNTDOWN_DURATION {
		return "0"
	} else {
		return fmt.Sprintf("%d", int(math.Ceil(timeRemaining.Minutes())))
	}
}

func (t *countdownTimer) Reset(app *tview.Application) {
	select {
	case t.reinitialiseSignal <- struct{}{}:
	default:
	}
	t.Start(app)
}

func initialiseRedraw(app *tview.Application, t *countdownTimer) time.Time {
	startTime := time.Now()
	app.QueueUpdateDraw(func() {
		t.countdown.SetText(t.startMessage)
	})
	time.Sleep(START_MESSAGE_DURATION)
	return startTime
}

func (t *countdownTimer) Start(app *tview.Application) {
	// noop if timer is running to ensure single redraw thread
	if !t.isRunning.CompareAndSwap(false, true) {
		return
	}

	select {
	case t.reinitialiseSignal <- struct{}{}:
	default:
	}

	go (func(app *tview.Application, t *countdownTimer) {
		startTime := initialiseRedraw(app, t)
		ticker := time.NewTicker(REFRESH_INTERVAL)
		defer ticker.Stop()

		for {
			select {
			case <-t.reinitialiseSignal:
				startTime = initialiseRedraw(app, t)
				continue
			case <-ticker.C:
			}

			minutesRemaining := getMinutesRemaining(startTime)
			if minutesRemaining == "0" {
				break
			}

			app.QueueUpdateDraw(func() {
				t.countdown.SetText(fmt.Sprintf("%s minutes remaining", minutesRemaining))
			})
		}

		app.QueueUpdateDraw(func() {
			t.countdown.SetText(t.endMessage)
		})

		t.isRunning.Store(false)
	})(app, t)
}
