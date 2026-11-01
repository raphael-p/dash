package database

import (
	"fmt"

	"github.com/raphael-p/dash/pkg/sqlutil"
	"github.com/raphael-p/dash/pkg/strutil"
)

type TaskInsertionOrder int

const (
	InsertFirst TaskInsertionOrder = iota + 1
	InsertNext
	InsertLast
)

type service struct {
	displayError   func(error)
	displayInfo    func(string)
	displayWarning func(string)
}

var instance service

func StartService(displayError func(error), displayInfo, displayWarning func(string)) {
	instance = service{displayError, displayInfo, displayWarning}
}

var Service = func() service {
	if instance.displayError != nil && instance.displayInfo != nil && instance.displayWarning != nil {
		return instance
	}
	panic("database service invoked before being started")
}

func (s service) AddTask(name, description string, insertionOrder TaskInsertionOrder, onSuccess func()) {
	if name == "" || description == "" {
		s.displayWarning("please provide a name and description of the task.")
		return
	}

	task, err := CreateTask(name, description)
	if err != nil {
		s.displayError(fmt.Errorf("%w: %w", errFailedCreate, err))
		return
	}

	switch insertionOrder {
	case InsertFirst:
		task.PriorityBumpedAt = sqlutil.NowNullTime()
		task.Update()
	case InsertNext:
		topTask, err := GetTopTask()
		if err != nil {
			s.displayError(fmt.Errorf("%w: %w", errFailedCreate, err))
			return
		}
		task.PriorityBumpedAt = sqlutil.NowNullTime()
		task.Update()
		topTask.PriorityBumpedAt = sqlutil.NowNullTime()
		topTask.Update()
	case InsertLast: // no action required
	default:
		s.displayError(errInvalidInsertion)
		return
	}

	s.displayInfo(fmt.Sprintf("new task [%d] created.", task.ID))
	onSuccess()
}

func (s service) GetTask(id int64) (Task, bool) {
	task, err := GetTask(int64(id))
	if err != nil {
		s.displayError(fmt.Errorf("%w [%d]: %w", errFailedGet, id, err))
		return task, false
	}
	return task, true
}

func (s service) GetTaskFromString(str string) *Task {
	if id, ok := strutil.ParseIDOrWarn(str, s.displayWarning); ok {
		if task, ok := s.GetTask(id); ok {
			return &task
		}
	}
	return nil
}

func (s service) RemoveTask(id int64) bool {
	t := Task{ID: id}
	deleted, err := t.Delete()
	if err != nil {
		s.displayError(fmt.Errorf("%w: %w", errFailedDelete, err))
		return false
	}

	if deleted {
		s.displayInfo(fmt.Sprintf("task [%d] deleted.", id))
	} else {
		s.displayWarning(fmt.Sprintf("task [%d] does not exist, noop.", id))
	}
	return deleted
}

func (s service) UpdateTask(task *Task, operation string) bool {
	updated, err := task.Update()
	if err != nil {
		s.displayError(fmt.Errorf("%w on task [%d]: %s: %w", errFailedUpdate, task.ID, operation, err))
		return false
	}

	if updated {
		s.displayInfo(fmt.Sprintf("task [%d]: %s", task.ID, operation))
	} else {
		s.displayWarning(fmt.Sprintf("task [%d] does not exist, noop", task.ID))
	}
	return updated
}
