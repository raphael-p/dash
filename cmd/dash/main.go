package main

import (
	"database/sql"
	"os"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/configreader"
	"github.com/raphael-p/datashard/pkg/logger"
)

type Config struct {
	DashDurationSeconds uint16 `json:"dash_duration_seconds"`
}

var config *Config = &Config{}

func main() {
	logger.Create(".", false)
	defer logger.Close()
	configreader.ReadConfigFile("dash", ".", config)

	db, err := sql.Open("sqlite3", "./datashard.db")
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
