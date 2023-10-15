package main

import "time"

const (
	EVENT_CUSTOMER_CREATED string = ".CUSTOMER.CREATED"
	EVENT_CUSTOMER_DONE    string = ".CUSTOMER.DONE"
)

const (
	CustomerStatusWaiting   StatusType = 0
	CustomerStatusActive    StatusType = 1
	CustomerStatusCompleted StatusType = 2
)

type Customer struct {
	Id        uint
	Status    StatusType
	Tasks     []Task
	TimeAdd   time.Time
	TimeLeave time.Time
}

func createCustomer(id uint, tasks []Task) *Customer {
	c := new(Customer)
	c.Id = id
	c.Tasks = tasks
	c.Status = CustomerStatusWaiting
	c.onCreated()
	return c
}

func (c *Customer) onCreated() {
	Announce(Event{EVENT_CUSTOMER_CREATED, c})
	c.Status = CustomerStatusActive
	c.TimeAdd = time.Now()
}

func (c *Customer) onDone() {
	Announce(Event{EVENT_CUSTOMER_DONE, c})
	c.Status = CustomerStatusCompleted
	c.TimeLeave = time.Now()
}
