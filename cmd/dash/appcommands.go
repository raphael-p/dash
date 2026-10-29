package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/raphael-p/dash/internal/database"
)

func handleAppCommand(command string) bool {
	switch command {
	case "extract":
		since, err := parseExtractCompleted(os.Args[2:])
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "extract: %v\n", err)
			os.Exit(2)
		}
		if err := extractCompleted(since); err != nil {
			fmt.Fprintf(os.Stderr, "extract: %v\n", err)
			os.Exit(1)
		}
	case "task":
		name, description, err := parseCreateTask(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: %v\n", err)
			os.Exit(2)
		}
		task, err := database.CreateTask(name, description)
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("new task created: %d\n", task.ID)
	default:
		return false
	}
	return true
}

func parseCreateTask(args []string) (string, string, error) {
	if len(args) != 2 {
		return "", "", errors.New("usage: dsh task <name> <description>")
	}
	name := args[0]
	description := args[1]
	if name == "" || description == "" {
		return "", "", errors.New("task name and description cannot be empty")
	}
	return name, description, nil
}

func parseExtractCompleted(args []string) (time.Time, error) {
	now := time.Now()
	since := now

	command := flag.NewFlagSet("extract", flag.ContinueOnError)
	command.SetOutput(os.Stderr)
	days := command.Int("days", 0, "include tasks completed in the last N days")
	weeks := command.Int("weeks", 0, "include tasks completed in the last N weeks")
	months := command.Int("months", 0, "include tasks completed in the last N months")
	years := command.Int("years", 0, "include tasks completed in the last N years")
	date := command.String("since", "", "include tasks completed since YYYY-MM-DD")
	if err := command.Parse(args); err != nil {
		return since, err
	}

	if command.NFlag() != 1 {
		return since, errors.New("exactly one of --days, --weeks, --months, --years, or --since is required")
	}
	var setFlagName string
	command.Visit(func(f *flag.Flag) { setFlagName = f.Name })

	switch setFlagName {
	case "since":
		parsed, err := time.Parse("2006-01-02", *date)
		if err != nil {
			return since, fmt.Errorf("--since must be in YYYY-MM-DD format: %w", err)
		}
		since = parsed
	case "days":
		if *days <= 0 {
			return since, errors.New("--days must be greater than zero")
		}
		since = now.AddDate(0, 0, -*days)
	case "weeks":
		if *weeks <= 0 {
			return since, errors.New("--weeks must be greater than zero")
		}
		since = now.AddDate(0, 0, -7*(*weeks))
	case "months":
		if *months <= 0 {
			return since, errors.New("--months must be greater than zero")
		}
		since = now.AddDate(0, -*months, 0)
	case "years":
		if *years <= 0 {
			return since, errors.New("--years must be greater than zero")
		}
		since = now.AddDate(-*years, 0, 0)
	}
	return since, nil
}

func extractCompleted(since time.Time) error {
	tasks, err := database.GetCompletedTasksSince(since)
	if err != nil {
		return err
	}
	rows := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		rows = append(rows, task.JSON())
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(rows)
}
