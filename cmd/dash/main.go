package main

import (
	"cmp"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/internal/tui/controller"
	"github.com/raphael-p/dash/pkg/logger"
	"github.com/rivo/tview"
)

var AppVersion = "dev"

func defaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".dash")
}

func main() {
	dataDir := cmp.Or(os.Getenv("DASH_DATA_DIR"), defaultDataDir())

	logger.Create(dataDir, false)
	defer logger.Close()

	config.Initialise(dataDir)

	db, err := sql.Open("sqlite3", filepath.Join(dataDir, "dash.db"))
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	database.DB = db

	if len(os.Args) < 2 {
		logger.Debugf("starting dash TUI (app version %s)", AppVersion)
		database.CheckVersion()
		app := tview.NewApplication()
		controller.Initialise(app)
		if err := app.Run(); err != nil {
			logger.Fatal(err.Error())
		}
		os.Exit(0)
	}

	command := os.Args[1]
	logger.Debugf("executing command %s (app version %s)", command, AppVersion)
	if handleSysCommand(command) || handleDBCommand(command) {
		os.Exit(0)
	}
	database.CheckVersion()
	if handleAppCommand(command) {
		os.Exit(0)
	}

	fmt.Printf("unknown command, usage: \n\n%s\n", helpText())
	os.Exit(2)
}
