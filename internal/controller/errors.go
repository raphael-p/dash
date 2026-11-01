package controller

import "errors"

var (
	errFailedCopy       = errors.New("failed to copy string")
	errFailedDashStart  = errors.New("failed to start dash")
	errFailedDashUpdate = errors.New("failed to update time spent on task")
	errFailedScrollDown = errors.New("failed to scroll down")
	errFailedScrollUp   = errors.New("failed to scroll up")
	errNoop             = errors.New("noop")
)
