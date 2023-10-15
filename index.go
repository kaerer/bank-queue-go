package main

import (
	"log"
	"strconv"
	"time"
)

const (
	EVENT_MANAGER_ALL string = "."
)

func createDemo(customerAmount int) []Customer {
	customers := make([]Customer, 0)
	for i := 0; i < customerAmount; i++ {
		var tasks = make([]Task, 0)
		idStr := strconv.Itoa(i)
		var taskAmount int = 1 //rand.Intn(5)
		for ii := 0; ii < taskAmount; ii++ {
			tasks = append(tasks, *createTask(idStr+"-"+strconv.Itoa(ii), "", 1)) // rand.Intn(5)
		}
		customer := *createCustomer(uint(i), tasks)
		customers = append(customers, customer)
	}
	return customers
}

func main() {

	customers := createDemo(10)
	m := *createManager(2, customers, 0)

	Verbosity = 1

	//- listen
	chn := Listen(EVENT_MANAGER_ALL)
	go func() {
		i := 0
		for e := range chn {
			log.Println(e.Tag, e.Data)

			// Avoid listening forever
			i++
			if m.Status != ManagerStatusActive {
				break
			}
		}
	}()

	m.init()
	m.start()

	time.Sleep(time.Duration(5 * time.Second))
}
