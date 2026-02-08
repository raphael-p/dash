package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/raphael-p/datashard/internal/database"
)

func handleCommand(command string) {
	switch command {
	case "init":
		database.Initialise()
	case "wipe":
		database.Wipe()
	case "use-sample-data":
		sampleCmd := flag.NewFlagSet("use-sample-data", flag.ExitOnError)
		randomEntryCount := sampleCmd.Int("randomEntryCount", 0, "number of random entries to generate on top of the default entries")
		sampleCmd.Parse(os.Args[2:])
		database.FillWithSampleData(*randomEntryCount)
	default:
		fmt.Println("unknown command, usage: datashard [init|use-sample-data|wipe|add|list|done|delete]")
		os.Exit(1)
	}
}
