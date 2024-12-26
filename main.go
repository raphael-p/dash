// main.go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/raphael-p/datashard/actions"
	"github.com/raphael-p/datashard/database"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: datashard [init|add|list]")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", "./datashard.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	database.DB = db

	command := os.Args[1]

	switch command {
	case "init":
		database.Initialize()
	case "wipe":
		fmt.Println("wipe not yet implemented")
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		title := addCmd.String("title", "", "Title of the note")
		content := addCmd.String("content", "", "Content of the note")
		urgency := addCmd.Int("urgency", 1, "Urgency (0-3)")
		dueDate := addCmd.String("due", "", "Due date (YYYY-MM-DD)")

		addCmd.Parse(os.Args[2:])

		if *title == "" || *content == "" {
			addCmd.Usage()
			os.Exit(1)
		}

		if *urgency < 1 || *urgency > 3 {
			fmt.Println("Urgency must be between 0 and 3.")
			os.Exit(1)
		}

		actions.AddNote(*title, *content, *urgency, *dueDate)
	case "list":
		actions.ListNotes()
	default:
		fmt.Println("Unknown command. Usage: datashard [init|add|list]")
		os.Exit(1)
	}
}
