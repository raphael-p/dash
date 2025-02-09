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
		name := addCmd.String("name", "", "name of the task")
		description := addCmd.String("description", "", "description of the task")

		addCmd.Parse(os.Args[2:])

		if *name == "" {
			addCmd.Usage()
			os.Exit(2)
		}

		actions.AddTask(*name, *description)
	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		query := listCmd.String("query", "", "search query to filter results on")
		showAll := listCmd.Bool("a", false, "show all tasks, including those already done")
		showDone := listCmd.Bool("d", false, "only show done tasks")

		listCmd.Parse(os.Args[2:])
		if *showAll && *showDone {
			fmt.Println("ambiguous command, do not use both -a and -d flags")
			listCmd.Usage()
			os.Exit(2)
		}

		var listMode actions.ListMode
		if *showAll {
			listMode = actions.ListAll
		} else if *showDone {
			listMode = actions.ListDone
		} else {
			listMode = actions.ListTodo
		}

		actions.ListTasks(listMode, *query)
	default:
		fmt.Println("unknown command, usage: datashard [init|wipe|add|list]")
		os.Exit(1)
	}
}
