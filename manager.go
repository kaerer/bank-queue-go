package main

import (
	"errors"
	"log"
	"strconv"
	"sync"
)

const (
	EVENT_MANAGER_INIT    string = ".MANAGER.INIT"
	EVENT_MANAGER_STARTED string = ".MANAGER.STARTED"
	EVENT_MANAGER_STOPPED string = ".MANAGER.STOPED"
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
	m.WorkerAmount = workerAmount
	m.Workers = make([]Worker, 0)

	// add workers
	for i := 0; i < workerAmount; i++ {
		worker := *createWorker(uint(i))
		worker.enable()
		m.Workers = append(m.Workers, worker)
	}
	return m
}

func (m *Manager) init() {
	Signal(EVENT_MANAGER_INIT)
	customers := m.createDemo(10)
	m.Queue.addMultipleCustomer(customers)
}

func (m *Manager) createDemo(customerAmount int) []Customer {
	customers := make([]Customer, 0)
	for i := 0; i < customerAmount; i++ {
		var tasks = make([]Task, 0)
		idStr := strconv.Itoa(i)
		tasks = append(tasks, *createTask(idStr+"-1", "", 1)) // rand.Intn(5)
		tasks = append(tasks, *createTask(idStr+"-2", "", 1)) // rand.Intn(5)
		customer := *createCustomer(uint(i), tasks)
		customers = append(customers, customer)
	}
	//TODO: generate customers with tasks for testing
	return customers
}

func (m Manager) getAvailableWorker() (*Worker, error) {
	if len(m.Workers) > 0 {
		for _, worker := range m.Workers {
			if worker.isAvailable() {
				log.Println("worker found")
				return &worker, nil
			}
		}
		log.Println("worker not found")
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
		}

		worker, err := m.getAvailableWorker()
		if err != nil {
			log.Println("Worker ERROR:")
			log.Println(err)
			break
		}
		if worker == nil {
			log.Println("no worker defined")
		}

		if worker == nil {
			continue
		}

		if customer == nil {
			m.stop()
		}

		m.wg.Add(1)
		go func(step int, worker Worker, customer Customer) {
			worker.setCustomer(&customer)
			worker.work()
			defer m.wg.Done()
			log.Println("'LIST'", customer.Id, worker.Id)
		}(step, *worker, *customer)

	}

	m.wg.Wait()
}

func (m *Manager) stop() {
	if len(m.Workers) > 0 {
		for _, worker := range m.Workers {
			worker.disable()
		}
	}

	m.Status = ManagerStatusPassive
}
