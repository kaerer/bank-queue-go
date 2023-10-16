package main

import (
	"errors"
)

const (
	EVENT_WORKER_CREATED        string = ".WORKER.CREATED"
	EVENT_WORKER_ENABLED        string = ".WORKER.ENABLED"
	EVENT_WORKER_DISABLED       string = ".WORKER.DISABLED"
	EVENT_WORKER_ASSIGNED       string = ".WORKER.ASSIGNED"
	EVENT_WORKER_TASK_STARTED   string = ".WORKER.TASK.STARTED"
	EVENT_WORKER_TASK_COMPLETED string = ".WORKER.TASK.COMPLETED"
)

const (
	WorkerStatusPassive  StatusType = 0
	WorkerStatusIdle     StatusType = 1
	WorkerStatusAssigned StatusType = 2
	WorkerStatusWorking  StatusType = 3
)

type Worker struct {
	Id              uint
	Status          StatusType
	CurrentCustomer *Customer
}

func createWorker(id uint) *Worker {
	w := new(Worker)
	w.Id = id
	w.CurrentCustomer = nil
	w.Status = WorkerStatusPassive
	Announce(Event{EVENT_WORKER_CREATED, w})
	return w
}

func (w *Worker) work() {
	w.onTasksStarted()
	var tasks []Task = w.CurrentCustomer.Tasks
	for _, task := range tasks {
		task.run()
	}
	w.CurrentCustomer.onDone()
	w.onTasksCompleted()
}

func (w *Worker) assignCustomer(customer *Customer) error {
	if w.CurrentCustomer != nil {
		return errors.New("CurrentCustomer is already assigned")
	}
	Announce(Event{EVENT_WORKER_ASSIGNED, []uint{w.Id, customer.Id}})
	w.Status = WorkerStatusAssigned
	w.CurrentCustomer = customer
	return nil
}

func (w Worker) isAvailable() bool {
	return w.Status == WorkerStatusIdle
}

func (w *Worker) enable() {
	Announce(Event{EVENT_WORKER_ENABLED, w.Id})
	w.Status = WorkerStatusIdle
}

func (w *Worker) disable() {
	Announce(Event{EVENT_WORKER_DISABLED, w.Id})
	w.Status = WorkerStatusPassive
}

func (w *Worker) onTasksStarted() {
	Announce(Event{EVENT_WORKER_TASK_STARTED, []uint{w.Id, w.CurrentCustomer.Id}})
	w.Status = WorkerStatusWorking
}

func (w *Worker) onTasksCompleted() {
	Announce(Event{EVENT_WORKER_TASK_COMPLETED, []uint{w.Id, w.CurrentCustomer.Id}})
	w.Status = WorkerStatusIdle
	w.CurrentCustomer = nil
}
