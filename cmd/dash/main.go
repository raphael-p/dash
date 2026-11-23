package main

import (
	"cmp"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/controller"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/cli"
	"github.com/raphael-p/dash/pkg/ghutil"
	"github.com/raphael-p/dash/pkg/logger"
	"github.com/raphael-p/dash/pkg/markerfile"
	"github.com/raphael-p/dash/pkg/sqlitebackup"
	"github.com/rivo/tview"
)

var AppVersion = "dev"
var dataDir = cmp.Or(os.Getenv("DASH_DATA_DIR"), defaultDataDir())

const gitRepo string = "raphael-p/dash"

func defaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".dash")
}

func installUpdate() {
	if err := ghutil.RunInstaller(gitRepo); err != nil {
		fmt.Printf("udpdate failed: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func proposeUpdate() {
	if markerfile.DebounceIsActive(dataDir, ".tagcheckdebounce", time.Hour*24) {
		return
	}
	latestVersion, err := ghutil.GetLatestTag(gitRepo)
	if err != nil {
		logger.Warningf("could not retrieve latest app version from github: %s", err)
		return
	}
	if latestVersion == AppVersion {
		return
	}
	if !cli.PromptConfirmation(fmt.Sprintf("dash update available (%s -> %s), install?", AppVersion, latestVersion)) {
		return
	}
	installUpdate()
}

func main() {
	// commands which don't require DB, config, or logger
	if len(os.Args) > 1 && handleSysCommand(os.Args[1]) {
		os.Exit(0)
	}
	proposeUpdate()

	// start the logger
	logger.Create(dataDir, false)
	defer logger.Close()

	// start the DB
	dbPath := filepath.Join(dataDir, "dash.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	database.DB = db

	// launch a conditional DB backup in a separate thread
	go func() {
		if err := sqlitebackup.Backup(db, dbPath); err != nil {
			logger.Warning(fmt.Sprintf("database backup failed: %v", err))
		}
	}()

	// start the TUI
	if len(os.Args) < 2 {
		config.Initialise(dataDir, tcell.ColorLimeGreen, tcell.ColorGray) // at the moment, config is only used by the TUI
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
	// operations on the database itself (init, wipe...)
	if handleDBCommand(command) {
		os.Exit(0)
	}

	// direct CLI interactions
	database.CheckVersion()
	if handleAppCommand(command) {
		os.Exit(0)
	}

	fmt.Printf("unknown command, usage: \n\n%s\n", helpText())
	os.Exit(2)
}
