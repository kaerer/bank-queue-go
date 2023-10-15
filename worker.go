package main

const (
	EVENT_WORKER_CREATED        string = ".WORKER.CREATED"
	EVENT_WORKER_ENABLED        string = ".WORKER.ENABLED"
	EVENT_WORKER_DISABLED       string = ".WORKER.DISABLED"
	EVENT_WORKER_TASK_STARTED   string = ".WORKER.TASK.STARTED"
	EVENT_WORKER_TASK_COMPLETED string = ".WORKER.TASK.COMPLETED"
)

const (
	WorkerStatusPassive StatusType = 0
	WorkerStatusIdle    StatusType = 1
	WorkerStatusWorking StatusType = 2
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
	w.Status = WorkerStatusIdle
	Announce(Event{EVENT_WORKER_CREATED, w})
	return w
}

func (w *Worker) work() {
	var tasks []Task = w.CurrentCustomer.Tasks
	w.onTasksStarted()
	for _, task := range tasks {
		task.run()
	}
	w.onTasksCompleted()
	w.CurrentCustomer.onDone()
}

func (w *Worker) setCustomer(customer *Customer) {
	w.CurrentCustomer = customer
}

func (w Worker) isAvailable() bool {
	return w.Status == WorkerStatusIdle
}

func (w *Worker) enable() {
	Announce(Event{EVENT_WORKER_ENABLED, &w})
	w.Status = WorkerStatusIdle
}

func (w *Worker) disable() {
	Announce(Event{EVENT_WORKER_DISABLED, &w})
	w.Status = WorkerStatusPassive
}

func (w *Worker) onTasksStarted() {
	Announce(Event{EVENT_WORKER_TASK_STARTED, nil})
	w.Status = WorkerStatusWorking
}

func (w *Worker) onTasksCompleted() {
	Announce(Event{EVENT_WORKER_TASK_COMPLETED, nil})
	w.Status = WorkerStatusIdle
}
