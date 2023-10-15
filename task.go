package main

import (
	"time"
)

const (
	TaskStatusRunning   StatusType = 2
	TaskStatusWaiting   StatusType = 1
	TaskStatusCompleted StatusType = 0
)

type Task struct {
	Status            StatusType
	Name              string
	Description       string
	WorkTimeInSeconds int
	TimeStarted       time.Time
	TimeCompleted     time.Time
}

func createTask(name string, description string, timeInSeconds int) *Task {
	t := new(Task)
	t.Name = name
	t.Description = description
	t.WorkTimeInSeconds = timeInSeconds
	t.Status = TaskStatusWaiting
	return t
}

func (t *Task) run() {
	t.onStarted()
	seconds := int64(t.WorkTimeInSeconds)
	time.Sleep(time.Duration(seconds) * time.Second)
	t.onCompleted()
}

func (t *Task) onStarted() {
	t.Status = TaskStatusRunning
	t.TimeStarted = time.Now()
}

func (t *Task) onCompleted() {
	t.Status = TaskStatusCompleted
	t.TimeCompleted = time.Now()
}
