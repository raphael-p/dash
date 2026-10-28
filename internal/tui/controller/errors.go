package controller

import "errors"

var (
	errInvalidInsertion = errors.New("invalid task insertion order")
	errFailedCopy       = errors.New("failed to copy string")
	errFailedCreate     = errors.New("failed to create task")
	errFailedDashStart  = errors.New("failed to start dash")
	errFailedDashUpdate = errors.New("failed to update time spent on task")
	errFailedDelete     = errors.New("failed to delete task")
	errFailedGet        = errors.New("failed to get task")
	errFailedMark       = errors.New("failed to mark task as done")
	errFailedUnmark     = errors.New("failed to unmark task as done")
	errFailedUpdate     = errors.New("failed update")
	errFailedScrollDown = errors.New("failed to scroll down")
	errFailedScrollUp   = errors.New("failed to scroll up")
	errNoop             = errors.New("noop")
)
