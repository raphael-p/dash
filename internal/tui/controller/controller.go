package controller

import (
	"time"

	"github.com/raphael-p/dash/internal/tui/panels"
	"github.com/raphael-p/dash/pkg/tviewcomponents/commandline"
	"github.com/rivo/tview"
)

type config struct {
	dashDuration         time.Duration
	descriptionCharLimit int
	nameCharLimit        int
}

type Controller struct {
	app          *tview.Application
	commandLine  *commandline.CommandLine
	navPanel     *panels.NavigationPanel
	infoPanel    *panels.InfoPanel
	displayPanel *panels.DisplayPanel
	config       config
}

func NewController(
	app *tview.Application,
	commandLine *commandline.CommandLine,
	navPanel *panels.NavigationPanel,
	infoPanel *panels.InfoPanel,
	displayPanel *panels.DisplayPanel,
	dashDurationSeconds uint16,
	descriptionCharLimit int,
	nameCharLimit int,
) *Controller {
	dashDuration := time.Duration(dashDurationSeconds) * time.Second
	config := config{dashDuration, descriptionCharLimit, nameCharLimit}
	return &Controller{app, commandLine, navPanel, infoPanel, displayPanel, config}
}
