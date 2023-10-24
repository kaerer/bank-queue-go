package main

import (
	"errors"
	"log"
	"sync"
)

const (
	EVENT_MANAGER_INIT               string = ".MANAGER.INIT"
	EVENT_MANAGER_STARTED            string = ".MANAGER.STARTED"
	EVENT_MANAGER_STOPPED            string = ".MANAGER.STOPED"
	EVENT_MANAGER_WORKER_FOUND       string = ".MANAGER.WORKER.FOUND"
	EVENT_MANAGER_CUSTOMER_NOT_FOUNT string = ".MANAGER.CUSTOMER.NOT.FOUND"
)

const (
	ManagerStatusActive  StatusType = 1
	ManagerStatusPassive StatusType = 0
)

type Manager struct {
	WorkerAmount int
	Workers      []*Worker
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
	m.Workers = make([]*Worker, 0)
	for i := 0; i < workerAmount; i++ {
		worker := createWorker(uint(i))
		worker.enable()
		m.Workers = append(m.Workers, worker)
	}
	return m
}

func (m *Manager) init() {
	Signal(EVENT_MANAGER_INIT)
}

func (m *Manager) getAvailableWorker() (*Worker, error) {
	if len(m.Workers) > 0 {
		for _, worker := range m.Workers {
			if worker.isAvailable() {
				Announce(Event{EVENT_MANAGER_WORKER_FOUND, worker.Id})
				return worker, nil
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
		// log.Println("Step: " + strconv.Itoa(step))

		if m.Status != ManagerStatusActive {
			break
		}

		step++

		worker, err := m.getAvailableWorker()
		if err != nil {
			log.Println("Worker ERROR:")
			log.Println(err)
			break
		}

		if worker == nil {
			continue
		}

		customer, _, err := m.Queue.getNextCustomer()
		if err != nil {
			log.Println("Customer ERROR:")
			log.Println(err)
			break
		}

		if customer == nil {
			Signal(EVENT_MANAGER_CUSTOMER_NOT_FOUNT)
			break
		}

		error := worker.assignCustomer(customer)
		if error == nil {
			m.wg.Add(1)
			go func(worker *Worker, customer *Customer, m *Manager) {
				worker.work()
				defer m.wg.Done()
				log.Println("Results: customer id: ", customer.Id, " customer group: ", customer.Group, " worker id: ", worker.Id)
			}(worker, customer, m)
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
	m.Status = ManagerStatusPassive
	if len(m.Workers) > 0 {
		for _, worker := range m.Workers {
			worker.disable()
		}
	}

	Signal(EVENT_MANAGER_STOPPED)
}
