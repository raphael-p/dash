package database

import (
	"database/sql"
	"strings"
)

func LazyInit(err error) bool {
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

func ExecWithLazyInit(tx DBExecutor, query string, args ...any) (sql.Result, error) {
	ex := getExecutor(tx)
	res, err := ex.Exec(query, args...)
	if LazyInit(err) {
		res, err = ex.Exec(query, args...)
	}
	return res, err
}
