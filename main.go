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

	command := os.Args[1]

	switch command {
	case "init":
		database.InitializeDatabase(db)
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		title := addCmd.String("title", "", "Title of the note")
		content := addCmd.String("content", "", "Content of the note")
		importance := addCmd.Int("importance", 1, "Importance level (1-3)")
		dueDate := addCmd.String("due", "", "Due date (YYYY-MM-DD)")

		addCmd.Parse(os.Args[2:])

		if *title == "" || *content == "" {
			addCmd.Usage()
			os.Exit(1)
		}

		if *importance < 1 || *importance > 3 {
			fmt.Println("Importance level must be between 1 and 3.")
			os.Exit(1)
		}

		actions.AddNote(db, *title, *content, *importance, *dueDate)
	case "list":
		actions.ListNotes(db)
	default:
		fmt.Println("Unknown command. Usage: cybernote [init|add|list]")
		os.Exit(1)
	}
}
