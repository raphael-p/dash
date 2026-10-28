package main

import (
	"fmt"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/cli"
)

func handleDBCommand(command string) bool {
	switch command {
	case "init":
		database.Initialise()
	case "wipe":
		if cli.PromptConfirmation("are you sure you want to wipe the database?") {
			database.Wipe()
		} else {
			fmt.Println("wipe cancelled")
		}
	default:
		return false
	}
	return true
}
