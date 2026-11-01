package database

import "errors"

var (
	ErrNotFound = errors.New("no task found")
)

var (
	errFailedCreate          = errors.New("failed to create task")
	errFailedDelete          = errors.New("failed to delete task")
	errFailedGet             = errors.New("failed to get task")
	errFailedUpdate          = errors.New("failed update")
	errInvalidInsertion      = errors.New("invalid task insertion order")
	errIDNotPositive         = errors.New("task ID must be greater than zero")
	errLastInsertFetchFailed = errors.New("failed to retrieve last insert ID")
)
