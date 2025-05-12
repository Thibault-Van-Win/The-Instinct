package scheduler

import "errors"

var (
	ErrTriggerNotFound = errors.New("trigger not found")
	ErrInvalidSchedule = errors.New("invalid schedule format")
	ErrAlreadyStarted  = errors.New("scheduler already started")
	ErrNotStarted      = errors.New("scheduler not started")
)
