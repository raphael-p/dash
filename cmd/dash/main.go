package main

import (
	"database/sql"
	"os"

	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/logger"
)

func main() {
	logger.Create(".", false)
	defer logger.Close()

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
