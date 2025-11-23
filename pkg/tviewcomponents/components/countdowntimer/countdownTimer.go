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

const (
	DEFAULT_COUNTDOWN_DURATION   = time.Minute * 20
	REFRESH_INTERVAL_FOR_MINUTES = time.Second * 2
	REFRESH_INTERVAL_FOR_SECONDS = time.Millisecond * 100
	START_MESSAGE_DURATION       = time.Second * 2
)

type timescale int

const (
	timescaleMinutes timescale = iota
	timescaleSeconds
)

type config struct {
	startMessage, endMessage string
	countdownDuration        time.Duration
	endAction                func()
}

type countdownTimer struct {
	Layout             *tview.Flex
	countdown          *tview.TextView
	description        *tview.TextView
	config             config
	isRunning          atomic.Bool
	reinitialiseSignal chan struct{}
}

var instance countdownTimer
var once sync.Once

var Instance = func(startMessage, endMessage string, countdownDuration time.Duration, endAction func()) *countdownTimer {
	config := config{
		startMessage, endMessage, countdownDuration, endAction,
	}
	once.Do(func() {
		countdown := tview.NewTextView()
		countdown.SetTextColor(tcell.ColorLimeGreen)
		countdown.SetBorderPadding(1, 0, 2, 2)

		description := tview.NewTextView().SetDynamicColors(true)
		description.SetBorderPadding(0, 1, 2, 2)

		layout := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(countdown, 3, 0, false).
			AddItem(description, 0, 1, false)

		if countdownDuration == 0 {
			countdownDuration = DEFAULT_COUNTDOWN_DURATION
		}

		instance = countdownTimer{
			Layout:             layout,
			countdown:          countdown,
			description:        description,
			isRunning:          atomic.Bool{},
			reinitialiseSignal: make(chan struct{}, 1),
		}
	})
	instance.config = config
	return &instance
}

func (t *countdownTimer) SetDescription(text string) {
	t.description.SetText(text)
}

func (t *countdownTimer) getTimeRemaining(startTime time.Time) (string, timescale, string) {
	timeElapsed := time.Since(startTime)
	timeRemaining := t.config.countdownDuration - timeElapsed
	if timeRemaining < 0 || timeElapsed < 0 {
		return "0", 0, ""
	}

	var timeRemainingRounded int
	var timescale timescale
	var unitName string
	if timeRemaining > time.Minute {
		timeRemainingRounded = int(math.Floor(timeRemaining.Minutes()))
		timescale = timescaleMinutes
		unitName = "minute"
	} else {
		timeRemainingRounded = int(math.Ceil(timeRemaining.Seconds()))
		timescale = timescaleSeconds
		unitName = "second"
	}
	if timeRemainingRounded > 1 {
		unitName = fmt.Sprint(unitName, "s")
	}
	return fmt.Sprint(timeRemainingRounded), timescale, unitName
}

func (t *countdownTimer) Reset(app *tview.Application) {
	select {
	case t.reinitialiseSignal <- struct{}{}:
	default:
	}
}

func initialiseRedraw(app *tview.Application, t *countdownTimer) time.Time {
	startTime := time.Now()
	app.QueueUpdateDraw(func() {
		t.countdown.SetText(t.config.startMessage)
	})
	time.Sleep(START_MESSAGE_DURATION)
	return startTime
}

func setTickerDuration(ticker *time.Ticker, timescale timescale) {
	if timescale == timescaleSeconds {
		ticker.Reset(REFRESH_INTERVAL_FOR_SECONDS)
	} else {
		ticker.Reset(REFRESH_INTERVAL_FOR_MINUTES)
	}
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

		ticker := time.NewTicker(1)
		defer ticker.Stop()

		_, initialTimescale, _ := t.getTimeRemaining(startTime)
		setTickerDuration(ticker, initialTimescale)

		for {
			select {
			case <-t.reinitialiseSignal:
				startTime = initialiseRedraw(app, t)
				continue
			case <-ticker.C:
			}

			timeRemaining, timescale, unitName := t.getTimeRemaining(startTime)
			if timeRemaining == "0" {
				break
			}

			if timescale != initialTimescale {
				setTickerDuration(ticker, timescale)
			}

			app.QueueUpdateDraw(func() {
				t.countdown.SetText(fmt.Sprintf("%s %s remaining", timeRemaining, unitName))
			})
		}

		app.QueueUpdateDraw(func() {
			t.countdown.SetText(t.config.endMessage)
			t.config.endAction()
		})

		t.isRunning.Store(false)
	})(app, t)
}
