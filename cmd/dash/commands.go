package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/raphael-p/datashard/internal/database"
)

func getConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/yes]: ", prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func handleCommand(command string) {
	switch command {
	case "init":
		database.Initialise()
	case "wipe":
		if getConfirmation("are you sure you want to wipe the database?") {
			database.Wipe()
		} else {
			fmt.Println("wipe cancelled")
		}
	case "generate":
		sampleCmd := flag.NewFlagSet("generate", flag.ExitOnError)
		randomEntryCount := sampleCmd.Int("randomEntryCount", 0, "number of random entries to generate on top of the default entries")
		sampleCmd.Parse(os.Args[2:])
		database.FillWithSampleData(*randomEntryCount)
	default:
		fmt.Println("unknown command, usage: datashard [init|generate|wipe|add|list|done|delete]")
		os.Exit(1)
	}
}
