package database

import "errors"

var (
	ErrIDNotPositive = errors.New("task ID must be greater than zero")
	ErrNotFound      = errors.New("no task found")
)

var (
	errLastInsertFetchFailed = errors.New("failed to retrieve last insert ID")
)
