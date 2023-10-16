package main

import "errors"

const (
	EVENT_QUEUE_CREATED        string = ".QUEUE.CREATED"
	EVENT_QUEUE_CUSTOMER_ADDED string = ".QUEUE.CUSTOMER.ADDED"
)

type Queue struct {
	CurrentCustomerIndex int
	Customers            []Customer
}

func createQueue(existingCustomers []Customer, startCustomerIndex int) *Queue {
	q := new(Queue)
	q.Customers = existingCustomers
	q.CurrentCustomerIndex = startCustomerIndex
	Signal(EVENT_QUEUE_CREATED)
	return q
}

func (q *Queue) getNextCustomer() (*Customer, int, error) {
	//TODO: fix and walk on slice
	if len(q.Customers) > 0 {
		if len(q.Customers) <= q.CurrentCustomerIndex {
			return nil, len(q.Customers), nil
		}
		var customer Customer = q.Customers[q.CurrentCustomerIndex]
		q.CurrentCustomerIndex++
		return &customer, len(q.Customers), nil
	}

	return nil, len(q.Customers), errors.New("customers are empty")
}

func (q *Queue) addCustomer(customer Customer) {
	Announce(Event{EVENT_QUEUE_CUSTOMER_ADDED, customer.Id})
	q.Customers = append(q.Customers, customer)
}
func (q *Queue) addMultipleCustomer(customers []Customer) {
	for _, customer := range customers {
		q.addCustomer(customer)
	}
}
