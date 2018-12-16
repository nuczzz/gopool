package gopool

import "sync"

type Goroutine interface {
	IsWorking() bool
}

type goroutine struct {
	//IsWorking status of the goroutine,
	//set IsWorking to true before handle request,
	//we can call the goroutine to handle new request if IsWorking == false
	isWorking bool

	//Lock lock of goroutine
	Lock sync.Mutex
}

func (g *goroutine) IsWorking() bool {
	return g.isWorking
}

func newGoroutine() Goroutine {
	return &goroutine{}
}
