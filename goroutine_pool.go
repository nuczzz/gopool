package gopool

import (
	"sync"
)

type goroutinePool struct {
	//coreGoroutineNum number of core goroutine.
	//we create a new goroutine to handle a new request
	//when goroutine number less than coreGoroutineNum,
	//even if some goroutines are not working.
	coreGoroutineNum int

	//maxGoroutineNum max number of goroutine.
	//maxGoroutineNum count all of the goroutines(include core goroutines).
	maxGoroutineNum int

	//lock lock of goroutine pool.
	lock sync.Mutex

	//index return goroutines[index] when get goroutine from goroutine pool, and index add 1
	index int

	//workGoroutine count of working goroutine
	workGoroutine int

	//goroutines array of goroutine
	goroutines []Goroutine
}

func (gp *goroutinePool) nextIndexAdd() {
}

func (gp *goroutinePool) getIndex() int {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	//Considering efficiency, DO NOT use
	//gp.index = (gp.index+1) % gp.maxGoroutineNum
	index := gp.index
	if gp.index == gp.maxGoroutineNum-1 {
		gp.index = 0
	} else {
		gp.index++
	}
	return index
}

func (gp *goroutinePool) Get() (Goroutine, error) {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	//todo
}
