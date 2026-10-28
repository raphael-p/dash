package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/raphael-p/dash/pkg/cli"
)

func handleSysCommand(command string) bool {
	switch command {
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
