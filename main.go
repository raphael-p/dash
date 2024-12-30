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
		fmt.Println("usage: datashard [init|add|list]")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", "./datashard.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	database.DB = db

	command := os.Args[1]

	switch command {
	case "init":
		database.Initialize()
	case "wipe":
		database.Wipe()
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		title := addCmd.String("title", "", "Title of the note")
		content := addCmd.String("content", "", "Content of the note")
		importance := addCmd.Int("importance", 0, "Importance (0-2)")
		dueDate := addCmd.String("due", "", "Due date (YYYY-MM-DD)")

		addCmd.Parse(os.Args[2:])

		if *title == "" || *content == "" {
			addCmd.Usage()
			os.Exit(1)
		}

		actions.AddNote(*title, *content, *importance, *dueDate)
	case "list":
		actions.ListNotes()
	default:
		fmt.Println("unknown command, usage: datashard [init|add|list]")
		os.Exit(1)
	}
}
