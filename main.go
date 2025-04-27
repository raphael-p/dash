// main.go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/raphael-p/datashard/actions"
	"github.com/raphael-p/datashard/database"
	"github.com/raphael-p/datashard/logger"
	"github.com/raphael-p/datashard/tui"
	"github.com/raphael-p/datashard/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: datashard [init|add|list]")
		os.Exit(1)
	}

	logger.Create(".", false)
	defer logger.Close()

	db, err := sql.Open("sqlite3", "./datashard.db")
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	database.DB = db

	command := os.Args[1]

	switch command {
	case "tui":
		tui.Home()
	case "init":
		database.Initialise()
	case "wipe":
		database.Wipe()
	case "use-sample-data":
		utils.FillDatabaseWithSampleData()
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
	case "done":
		doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
		id := doneCmd.Int64("id", 0, "id of task to mark as done")

		doneCmd.Parse(os.Args[2:])
		if *id == 0 {
			doneCmd.Usage()
			os.Exit(2)
		}
		actions.MarkTaskAsDone(*id)
	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		id := deleteCmd.Int64("id", 0, "id of task to delete")

		deleteCmd.Parse(os.Args[2:])
		if *id == 0 {
			deleteCmd.Usage()
			os.Exit(2)
		}
		actions.DeleteTask(*id)
	default:
		fmt.Println("unknown command, usage: datashard [tui|init|use-sample-data|wipe|add|list|done|delete]")
		os.Exit(1)
	}
}
