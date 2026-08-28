package main

import (
	"cmp"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/configreader"
	"github.com/raphael-p/datashard/pkg/logger"
)

var AppVersion = "dev"

type Config struct {
	DashDurationSeconds  uint16 `json:"dash_duration_seconds"`
	DescriptionCharLimit int    `json:"description_char_limit"`
	NameCharLimit        int    `json:"name_char_limit"`
}

var config *Config = &Config{}

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

	configreader.ReadConfigFile("dash", dataDir, config)

	db, err := sql.Open("sqlite3", filepath.Join(dataDir, "datashard.db"))
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	database.DB = db

	database.CheckVersion()

	if len(os.Args) < 2 {
		logger.Debugf("starting dash TUI (app version %s)", AppVersion)
		startTUI()
	} else {
		command := os.Args[1]
		logger.Debugf("executing command %s (app version %s)", command, AppVersion)
		handleCommand(command)
	}
}
