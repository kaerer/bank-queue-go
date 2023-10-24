package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// - Settings
const workerAmount int = 15
const initialCustomerAmount int = 150

const maxTaskAmount int = 5
const maxTaskWorkTimeInSeconds int = 3

const newCustomerPeriodInSeconds int = 5
const newCustomerAmountInPeriod int = 1

// - ----

const (
	EVENT_ALL string = "."
)

func getRandomCount(max int) int {
	return rand.Intn(max-1) + 1
}

func createDemo(customerAmount int, customerGroup CustomerGroupType) []Customer {
	customers := make([]Customer, 0)
	for i := 0; i < customerAmount; i++ {
		var tasks = make([]Task, 0)
		idStr := strconv.Itoa(i)
		var taskAmount int = getRandomCount(maxTaskAmount)
		for ii := 0; ii < taskAmount; ii++ {
			tasks = append(tasks, *createTask(idStr+"-"+strconv.Itoa(ii), "", getRandomCount(maxTaskWorkTimeInSeconds)))
		}
		customer := *createCustomer(customerGroup, tasks)
		customers = append(customers, customer)
	}
	return customers
}

var m Manager
var wg sync.WaitGroup

func main() {
	Verbosity = 0
	customers := createDemo(initialCustomerAmount, CustomerGroupA)
	m = *createManager(workerAmount, customers, 0)

	if newCustomerPeriodInSeconds > 0 {
		wg.Add(1)
		go func() {
			for range time.Tick(time.Second * time.Duration(newCustomerPeriodInSeconds)) {
				if m.Status != ManagerStatusActive {
					log.Println("Doors are CLOSED")
					defer wg.Done()
					break
				}

				log.Println("New Customer Entered")
				newCustomer()
			}
		}()
	}

	wg.Add(1)
	go func() {
		// listen all event
		chn := Listen(EVENT_ALL)
		wg.Add(1)
		go func() {
			for e := range chn {
				log.Println(e.Tag, e.Data)

				// Avoid listening forever
				if m.Status != ManagerStatusActive {
					log.Println("Event Listen ENDS")
					break
				}
			}
			defer wg.Done()
		}()

		m.init()
		m.start()
		defer func() {
			wg.Done()
			log.Println("Main DONE")
		}()
	}()

	wg.Wait()

	time.Sleep(time.Duration(5 * time.Second))
}

func newCustomer() {
	customers := createDemo(newCustomerAmountInPeriod, CustomerGroupB)
	m.Queue.addMultipleCustomer(customers)
}
