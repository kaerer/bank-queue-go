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

const (
	CustomerGroupA CustomerGroupType = 0
	CustomerGroupB CustomerGroupType = 1
)

type Customer struct {
	Id        uint
	Status    StatusType
	Group     CustomerGroupType
	Tasks     []Task
	TimeAdd   time.Time
	TimeLeave time.Time
}

var maxCustomerId uint = 0

func createCustomer(group CustomerGroupType, tasks []Task) *Customer {
	maxCustomerId++

	c := new(Customer)
	c.Id = maxCustomerId
	c.Tasks = tasks
	c.Status = CustomerStatusWaiting
	c.onCreated()
	return c
}

func (c *Customer) onCreated() {
	Announce(Event{EVENT_CUSTOMER_CREATED, c.Id})
	c.Status = CustomerStatusActive
	c.TimeAdd = time.Now()
}

func (c *Customer) onDone() {
	Announce(Event{EVENT_CUSTOMER_DONE, c.Id})
	c.Status = CustomerStatusCompleted
	c.TimeLeave = time.Now()
}
