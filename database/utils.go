package database

import (
	"database/sql"
	"log"
	"strings"
	"sync"
)

func lazyInit(err error) bool {
	if err != nil && strings.Contains(err.Error(), "no such table") {
		Initialize()
		return true
	}
	return false
}

func QueryWithLazyInit(query string, args ...any) (*sql.Rows, error) {
	res, err := DB.Query(query, args...)
	if lazyInit(err) {
		res, err = DB.Query(query, args...)
	}
	return res, err
}

func ExecWithLazyInit(tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	res, err := tx.Exec(query, args...)
	if lazyInit(err) {
		res, err = tx.Exec(query, args...)
	}
	return res, err
}

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
