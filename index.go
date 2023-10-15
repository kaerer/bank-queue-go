package main

import (
	"fmt"
)

const (
	EVENT_MANAGER_ALL string = "."
)

func main() {

	m := *createManager(2, []Customer{}, 0)

	Verbosity = 1

	//- listen
	chn := Listen(EVENT_MANAGER_ALL)
	go func() {
		i := 0
		for e := range chn {
			fmt.Println(e.Tag, e.Data)

			// Avoid listening forever
			i++
			if m.Status != ManagerStatusActive {
				break
			}
		}
	}()

	m.init()
	m.start()

	// time.Sleep(time.Duration(20 * time.Second))
}
