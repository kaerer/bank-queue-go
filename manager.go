package main

import (
	"errors"
	"log"
	"strconv"
	"sync"
)

const (
	EVENT_MANAGER_INIT         string = ".MANAGER.INIT"
	EVENT_MANAGER_STARTED      string = ".MANAGER.STARTED"
	EVENT_MANAGER_STOPPED      string = ".MANAGER.STOPED"
	EVENT_MANAGER_WORKER_FOUND string = ".MANAGER.WORKER.FOUND"
)

const (
	ManagerStatusActive  StatusType = 1
	ManagerStatusPassive StatusType = 0
)

type Manager struct {
	WorkerAmount int
	Workers      []Worker
	Queue        Queue
	Status       StatusType
	wg           sync.WaitGroup
}

func createManager(workerAmount int, existingCustomers []Customer, startCustomerIndex int) *Manager {
	m := new(Manager)
	m.Queue = *createQueue(make([]Customer, 0), 0)
	m.Queue.CurrentCustomerIndex = startCustomerIndex
	m.WorkerAmount = workerAmount
	if existingCustomers != nil {
		m.Queue.addMultipleCustomer(existingCustomers)
	}

	// add workers
	m.Workers = make([]Worker, 0)
	for i := 0; i < workerAmount; i++ {
		log.Println(i)

		worker := *createWorker(uint(i))
		worker.enable()
		m.Workers = append(m.Workers, worker)
	}
	return m
}

func (m *Manager) init() {
	Signal(EVENT_MANAGER_INIT)
	// customers := m.createDemo(10)
	// m.Queue.addMultipleCustomer(customers)
}

// var previousIndex int = 0

func (m *Manager) getAvailableWorker() (*Worker, error) {
	// var worker Worker
	// if previousIndex == 0 {
	// 	previousIndex = 1
	// } else {
	// 	previousIndex = 0
	// }
	// worker = m.Workers[previousIndex]
	// Announce(Event{EVENT_MANAGER_WORKER_FOUND, worker.Id})
	// return &worker, nil

	if len(m.Workers) > 0 {
		for _, worker := range m.Workers {
			if worker.isAvailable() {
				Announce(Event{EVENT_MANAGER_WORKER_FOUND, worker.Id})
				return &worker, nil
			}
		}
		return nil, nil
	}

	return nil, errors.New("no available worker")
}

func (m *Manager) start() {

	m.Status = ManagerStatusActive
	Signal(EVENT_MANAGER_STARTED)

	// check if any waiting user in the queue, if not wait for more? or just stop
	// wait for any worker to be free back if there is any

	var step int = 0
	for {
		log.Println("Step: " + strconv.Itoa(step))

		if m.Status != ManagerStatusActive {
			break
		}

		step++

		customer, customerAmount, err := m.Queue.getNextCustomer()
		if err != nil {
			log.Println("Customer ERROR:")
			log.Println(err)
			break
		}

		if customer == nil {
			if customerAmount > 0 {
				log.Println("no customer left")
			} else {
				log.Println("no customer added")
			}

			break
		}

		worker, err := m.getAvailableWorker()
		if err != nil {
			log.Println("Worker ERROR:")
			log.Println(err)
			break
		}

		if worker == nil {
			log.Println("no worker defined")
			continue
		}

		error := worker.assignCustomer(customer)

		if error == nil {
			m.wg.Add(1)
			go func() {
				worker.work()
				defer m.wg.Done()
				log.Println("'Results:' customer id: ", customer.Id, " worker id: ", worker.Id)
			}()
		} else {
			log.Println("Worker ERROR:")
			log.Println(error)
			break
		}

	}

	m.wg.Wait()
	m.stop()
}

func (m *Manager) stop() {
	if len(m.Workers) > 0 {
		for _, worker := range m.Workers {
			worker.disable()
		}
	}

	m.Status = ManagerStatusPassive
}
