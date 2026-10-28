package database

import "errors"

var (
	ErrNotFound = errors.New("no task found")
)

var (
	errIDNotPositive         = errors.New("task ID must be greater than zero")
	errLastInsertFetchFailed = errors.New("failed to retrieve last insert ID")
)
