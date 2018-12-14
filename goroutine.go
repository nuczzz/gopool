package gopool

import "sync"

type Goroutine struct {
	//IsWorking status of the goroutine,
	//set IsWorking to true before handle request,
	//we can call the goroutine to handle new request if IsWorking == false
	IsWorking bool

	//Lock lock of goroutine
	Lock sync.Mutex
}
