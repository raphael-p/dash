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

func getExecutor(tx DBExecutor) DBExecutor {
	if tx != nil {
		return tx
	}
	return DB
}

type DBExecutor interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

func QueryWithLazyInit(tx DBExecutor, query string, args ...any) (*sql.Rows, error) {
	ex := getExecutor(tx)
	res, err := ex.Query(query, args...)
	if lazyInit(err) {
		res, err = ex.Query(query, args...)
	}
	return res, err
}

func ExecWithLazyInit(tx DBExecutor, query string, args ...any) (sql.Result, error) {
	ex := getExecutor(tx)
	res, err := ex.Exec(query, args...)
	if lazyInit(err) {
		res, err = ex.Exec(query, args...)
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
