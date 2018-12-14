package gopool

import "sync"

type GoroutinePool struct {
	//coreGoroutineNum number of core goroutine.
	//we create a new goroutine to handle a new request
	//when goroutine number less than coreGoroutineNum,
	//even if some goroutines are not working.
	coreGoroutineNum int

	//maxGoroutineNum max number of goroutine.
	//maxGoroutineNum count all of the goroutines(include core goroutines).
	maxGoroutineNum int

	//Lock lock of goroutine pool.
	Lock sync.Mutex
}
