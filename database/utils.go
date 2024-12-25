package database

import (
	"database/sql"
	"log"
	"sync"
)

func execAsync(db *sql.DB, command, errMessage string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := db.Exec(command)
		if err != nil {
			log.Fatalf("%s: %v", errMessage, err)
		}
	}()
}
