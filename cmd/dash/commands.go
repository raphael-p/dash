package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/raphael-p/dash/internal/database"
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
	case "extract":
		if err := extractCompleted(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "extract: %v\n", err)
			os.Exit(2)
		}
	default:
		fmt.Println("unknown command, usage: dsh [init|generate|wipe|extract]")
		os.Exit(1)
	}
}

func extractCompleted(args []string) error {
	command := flag.NewFlagSet("extract", flag.ContinueOnError)
	command.SetOutput(os.Stderr)
	days := command.Int("days", 0, "include tasks completed in the last N days")
	weeks := command.Int("weeks", 0, "include tasks completed in the last N weeks")
	months := command.Int("months", 0, "include tasks completed in the last N months")
	years := command.Int("years", 0, "include tasks completed in the last N years")
	since := command.String("since", "", "include tasks completed since YYYY-MM-DD")
	if err := command.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	if command.NFlag() != 1 {
		return errors.New("exactly one of --days, --weeks, --months, --years, or --since is required")
	}
	var setFlagName string
	command.Visit(func(f *flag.Flag) { setFlagName = f.Name })

	now := time.Now()
	start := now
	switch setFlagName {
	case "since":
		parsed, err := time.Parse("2006-01-02", *since)
		if err != nil {
			return fmt.Errorf("--since must be in YYYY-MM-DD format: %w", err)
		}
		start = parsed
	case "days":
		if *days <= 0 {
			return errors.New("--days must be greater than zero")
		}
		start = now.AddDate(0, 0, -*days)
	case "weeks":
		if *weeks <= 0 {
			return errors.New("--weeks must be greater than zero")
		}
		start = now.AddDate(0, 0, -7*(*weeks))
	case "months":
		if *months <= 0 {
			return errors.New("--months must be greater than zero")
		}
		start = now.AddDate(0, -*months, 0)
	case "years":
		if *years <= 0 {
			return errors.New("--years must be greater than zero")
		}
		start = now.AddDate(-*years, 0, 0)
	}
	tasks, err := database.GetCompletedTasksSince(start)
	if err != nil {
		return err
	}
	rows := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		rows = append(rows, taskAsJSON(task))
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(rows)
}

func taskAsJSON(task database.Task) map[string]any {
	return map[string]any{
		"id": task.ID, "name": task.Name, "description": task.Description,
		"created_at": task.CreatedAt, "updated_at": task.UpdatedAt,
		"completed_at":       nullableTime(task.CompletedAt),
		"priority_bumped_at": nullableTime(task.PriorityBumpedAt),
		"time_spent_seconds": nullableInt(task.TimeSpentSeconds),
	}
}

func nullableTime(value sql.NullTime) any {
	if value.Valid {
		return value.Time
	}
	return nil
}

func nullableInt(value sql.NullInt16) any {
	if value.Valid {
		return value.Int16
	}
	return nil
}
