package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/cli"
)

func handleDBCommand(command string) bool {
	switch command {
	case "init":
		database.Initialise()
	case "db:reset":
		reset := flag.NewFlagSet("reset", flag.ExitOnError)
		generateCount := reset.Int("generate", 0, "number of random entries par tab to generate on top of the default entries")
		err := reset.Parse(os.Args[2:])
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "db:reset:dev: %v\n", err)
			os.Exit(2)
		}
		if !cli.PromptConfirmation("are you sure you want to wipe the database and fill it with sample data?") {
			fmt.Println("reset cancelled")
			return true
		}
		database.FillWithSampleData(*generateCount)
	case "db:wipe":
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
