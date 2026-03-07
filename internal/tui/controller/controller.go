package controller

import (
	"time"

	"github.com/raphael-p/datashard/internal/tui/panels"
	"github.com/rivo/tview"
)

type config struct {
	dashDuration         time.Duration
	descriptionCharLimit int
	nameCharLimit        int
}

type Controller struct {
	app          *tview.Application
	inputPanel   *panels.InputPanel
	infoPanel    *panels.InfoPanel
	displayPanel *panels.DisplayPanel
	config       config
}

func NewController(
	app *tview.Application,
	inputPanel *panels.InputPanel,
	infoPanel *panels.InfoPanel,
	displayPanel *panels.DisplayPanel,
	dashDurationSeconds uint16,
	descriptionCharLimit int,
	nameCharLimit int,
) *Controller {
	dashDuration := time.Duration(dashDurationSeconds) * time.Second
	config := config{dashDuration, descriptionCharLimit, nameCharLimit}
	return &Controller{app, inputPanel, infoPanel, displayPanel, config}
}
