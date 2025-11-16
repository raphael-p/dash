package controller

import (
	"time"

	"github.com/raphael-p/datashard/internal/tui/panels"
	"github.com/rivo/tview"
)

type Controller struct {
	app          *tview.Application
	inputPanel   *panels.InputPanel
	infoPanel    *panels.InfoPanel
	displayPanel *panels.DisplayPanel
	dashDuration time.Duration
}

func NewController(
	app *tview.Application,
	inputPanel *panels.InputPanel,
	infoPanel *panels.InfoPanel,
	displayPanel *panels.DisplayPanel,
	dashDuration time.Duration,
) *Controller {
	return &Controller{app, inputPanel, infoPanel, displayPanel, dashDuration}
}
