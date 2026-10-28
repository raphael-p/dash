package countdowntimer

import (
	"cmp"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	DefaultCountdownDuration = time.Minute * 25
	DefaultSideEffectPeriod  = time.Second * 30
	LongRefreshInterval      = time.Second * 2
	ShortRefreshInterval     = time.Millisecond * 100
	StartMessageDuration     = time.Second * 2
)

type timescale int

const (
	timescaleMinutes timescale = iota
	timescaleSeconds
)

type Config struct {
	StartMessage, EndMessage string
	CountdownDuration        time.Duration
	PeriodicSideEffect       func(lastRunTime time.Time, isEnd bool) // function that runs periodically and on timer end
	SideEffectPeriod         time.Duration
}

type CountdownTimer struct {
	Layout             *tview.Flex
	countdown          *tview.TextView
	text               *tview.TextView
	config             Config
	isRunning          atomic.Bool
	reinitialiseSignal chan struct{}
	interruptSignal    chan struct{}
}

var instance CountdownTimer
var once sync.Once

var Instance = func() (*CountdownTimer, bool) {
	isNew := false
	once.Do(func() {
		countdown := tview.NewTextView()
		countdown.SetTextColor(tcell.ColorLimeGreen)

		text := tview.NewTextView().SetDynamicColors(true)

		layout := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(countdown, 3, 0, false).
			AddItem(text, 0, 1, false)

		instance = CountdownTimer{
			Layout:             layout,
			countdown:          countdown,
			text:               text,
			isRunning:          atomic.Bool{},
			reinitialiseSignal: make(chan struct{}, 1),
			interruptSignal:    make(chan struct{}, 1),
		}
		isNew = true
	})
	return &instance, isNew
}

func (t *CountdownTimer) SetConfig(config Config) {
	config.SideEffectPeriod = cmp.Or(config.SideEffectPeriod, DefaultSideEffectPeriod)
	config.CountdownDuration = cmp.Or(config.CountdownDuration, DefaultCountdownDuration)
	instance.config = config
}

func (t *CountdownTimer) SetInputCapture(ip func(event *tcell.EventKey) *tcell.EventKey) *tview.Box {
	return t.Layout.SetInputCapture(ip)
}

func (t *CountdownTimer) SetText(text string) *tview.TextView {
	return t.text.SetText(text)
}

func (t *CountdownTimer) getTimeRemaining(startTime time.Time) (string, timescale, string) {
	timeElapsed := time.Since(startTime)
	timeRemaining := t.config.CountdownDuration - timeElapsed
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

func (t *CountdownTimer) Reset() {
	select {
	case t.reinitialiseSignal <- struct{}{}:
	default:
	}
}

func (t *CountdownTimer) Stop() {
	select {
	case t.interruptSignal <- struct{}{}:
	default:
	}
}

func initialiseRedraw(queueUpdateDraw func(func()), t *CountdownTimer) time.Time {
	startTime := time.Now()
	queueUpdateDraw(func() { t.countdown.SetText(t.config.StartMessage) })
	select {
	case <-time.After(StartMessageDuration):
	case <-t.interruptSignal:
		t.Stop() // re-emit signal
	}
	return startTime
}

func setTickerDuration(ticker *time.Ticker, timescale timescale) {
	if timescale == timescaleSeconds {
		ticker.Reset(ShortRefreshInterval)
	} else {
		ticker.Reset(LongRefreshInterval)
	}
}

func (t *CountdownTimer) Start(queueUpdateDraw func(func())) {
	// noop if timer is running to ensure single redraw thread
	if !t.isRunning.CompareAndSwap(false, true) {
		return
	}

	// Drain any leftover interrupt signal before starting
	select {
	case <-t.interruptSignal:
	default:
	}

	select {
	case t.reinitialiseSignal <- struct{}{}:
	default:
	}

	go (func(queueUpdateDraw func(func()), t *CountdownTimer) {
		defer func() {
			t.isRunning.Store(false)
			queueUpdateDraw(func() { t.countdown.SetText(t.config.EndMessage) })
		}()

		startTime := initialiseRedraw(queueUpdateDraw, t)
		lastSideEffectRun := startTime

		ticker := time.NewTicker(1)
		defer ticker.Stop()

		sideEffectTicker := time.NewTicker(t.config.SideEffectPeriod)
		defer sideEffectTicker.Stop()

		_, initialTimescale, _ := t.getTimeRemaining(startTime)
		setTickerDuration(ticker, initialTimescale)

		for {
			select {
			case <-t.interruptSignal:
				return
			case <-t.reinitialiseSignal:
				startTime = initialiseRedraw(queueUpdateDraw, t)
				sideEffectTicker.Reset(t.config.SideEffectPeriod)
				lastSideEffectRun = startTime
				continue
			case <-sideEffectTicker.C:
				now := time.Now()
				t.config.PeriodicSideEffect(lastSideEffectRun, false)
				lastSideEffectRun = now
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

			queueUpdateDraw(func() { t.countdown.SetText(fmt.Sprintf("%s %s remaining", timeRemaining, unitName)) })
		}

		t.config.PeriodicSideEffect(lastSideEffectRun, true)
	})(queueUpdateDraw, t)
}
