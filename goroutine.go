package gopool

import (
	"sync"
	"sync/atomic"
)

type Goroutine interface {
	Execute(task func())
	IsWorking() bool
	Run()
}

type goroutine struct {
	//IsWorking status of the goroutine,
	//set IsWorking to true before handle request,
	//we can call the goroutine to handle new request if IsWorking == false
	isWorking int32

	//Lock lock of goroutine
	Lock sync.Mutex

	task chan func()
}

func (g *goroutine) IsWorking() bool {
	return atomic.LoadInt32(&g.isWorking) != 0
}

func (g *goroutine) Execute(task func()) {
	g.task <- task
}

func (g *goroutine) Run() {
	NewGoroutineWithRecover(func() {
		for f := range g.task {
			if f == nil {
				//kill goroutine with nil task
				return
			}

			f()
		}
	})
}

func newGoroutine() Goroutine {
	return &goroutine{
		task: make(chan func(), 1),
	}
}
