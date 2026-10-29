package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/raphael-p/dash/pkg/cli"
)

func handleSysCommand(command string) bool {
	switch command {
	case "--help", "--h", "-h":
		fallthrough
	case "help":
		fmt.Print(helpText())
	case "update":
		if !cli.PromptConfirmation("are you sure you want to update dash?") {
			fmt.Println("update cancelled")
			os.Exit(0)
		}
		cmdStr := "curl -fsSL https://github.com/raphael-p/dash/releases/latest/download/install.sh | sh"
		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("udpdate failed: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println(AppVersion)
	default:
		return false
	}
	return true
}

func helpText() string {
	instructions := []string{
		"dsh [-h]                          run dash",
		"dsh extract [flags] [-h]          output a JSON of tasks completed",
		"dsh help                          get CLI instructions",
		"dsh init                          initialise the database",
		"dsh task <name> <description>     create a new task",
		"dsh update                        fetch and install the latest version of dash from github",
		"dsh version                       print the current version of dash",
		"",
		"dsh db:reset [--generate] [-h]    wipes the database and populates it with sample data",
		"dsh db:wipe                       deletes all data",
	}
	return strings.Join(instructions, "\n") + "\n"
}
