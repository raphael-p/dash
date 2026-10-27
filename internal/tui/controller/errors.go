package controller

import "errors"

var (
	errEmptyID          = errors.New("task ID is empty")
	errIDNotInt         = errors.New("task ID is not an integer")
	errInvalidInsertion = errors.New("invalid task insertion order")
	errFailedBump       = errors.New("failed to bump task priority")
	errFailedCopy       = errors.New("failed to copy string")
	errFailedCreate     = errors.New("failed to create task")
	errFailedDashStart  = errors.New("failed to start dash")
	errFailedDashUpdate = errors.New("failed to update time spent on task")
	errFailedDelete     = errors.New("failed to delete task")
	errFailedEdit       = errors.New("failed to edit task")
	errFailedMark       = errors.New("failed to mark task as done")
	errFailedUnmark     = errors.New("failed to unmark task as done")
	errFailedPause      = errors.New("failed to pause task")
	errFailedUnpause    = errors.New("failed to unpause task")
	errFailedReactivate = errors.New("failed to reactivate task")
	errFailedScrollDown = errors.New("failed to scroll down")
	errFailedScrollUp   = errors.New("failed to scroll up")
	errNoop             = errors.New("noop")
)
