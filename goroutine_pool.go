package gopool

import "sync"

type GoroutinePool struct {
	//coreGoroutineNum the number of core goroutine.
	//we create a new goroutine to handle the new request
	//when goroutine number less than coreGoroutineNum,
	//even if some goroutines are not working.
	coreGoroutineNum int

	//maxGoroutineNum the max number of goroutine.
	//maxGoroutineNum count all of the goroutines(include core goroutine).
	maxGoroutineNum int

	//IsWorking status of the goroutine,
	//set IsWorking to true before handle request,
	//we can call the goroutine to handle new request if IsWorking == false
	IsWorking bool

	//Lock lock of the goroutine.
	Lock sync.Mutex
}
