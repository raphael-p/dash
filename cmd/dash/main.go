package main

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/configreader"
	"github.com/raphael-p/datashard/pkg/logger"
)

type Config struct {
	DashDurationSeconds  uint16 `json:"dash_duration_seconds"`
	DescriptionCharLimit int    `json:"description_char_limit"`
	NameCharLimit        int    `json:"name_char_limit"`
}

var config *Config = &Config{}

func main() {
	dataDir := os.Getenv("DASH_DATA_DIR")
	if dataDir == "" {
		dataDir = "."
	}
	logger.Create(dataDir, false)
	defer logger.Close()
	configreader.ReadConfigFile("dash", dataDir, config)

	db, err := sql.Open("sqlite3", filepath.Join(dataDir, "datashard.db"))
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	database.DB = db

	if len(os.Args) < 2 {
		startTUI()
	} else {
		handleCommand(os.Args[1])
	}
}
