package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/raphael-p/datashard/internal/cli/actions"
	"github.com/raphael-p/datashard/internal/database"
)

func handleCommand(command string) {
	switch command {
	case "init":
		database.Initialise()
	case "wipe":
		database.Wipe()
	case "use-sample-data":
		sampleCmd := flag.NewFlagSet("use-sample-data", flag.ExitOnError)
		randomEntryCount := sampleCmd.Int("randomEntryCount", 0, "number of random entries to generate on top of the default entries")
		sampleCmd.Parse(os.Args[2:])
		database.FillWithSampleData(*randomEntryCount)
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
		fmt.Println("unknown command, usage: datashard [init|use-sample-data|wipe|add|list|done|delete]")
		os.Exit(1)
	}
}
