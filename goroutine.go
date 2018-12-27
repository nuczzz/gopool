package gopool

import (
	"bytes"
	"log"
	"runtime/debug"
)

type Task func()

type Goroutine interface {
	Execute(task Task)
}

type goroutine struct {
	// pool who manage this goroutine.
	pool *goroutinePool

	// task channel for send or receive task.
	task chan Task
}

// Execute stop the timer before a goroutine execute task,
// and returns error if stop timer failed(timer has stopped).
// When g.timer received timeout signal, goroutine pool will
// release the goroutine, and the task will not be executed
// forever.
func (g *goroutine) Execute(task Task) {
	g.task <- task
}

func (g *goroutine) run() {
	g.newGoroutineWithRecover(func() {
		for {
			select {
			case task := <-g.task:
				task()
				g.pool.recycleGoroutine(g)
			}
		}
	})
}

func newGoroutine(pool *goroutinePool) Goroutine {
	g := &goroutine{
		pool: pool,
		task: make(chan Task),
	}
	g.run()
	return g
}

func (g *goroutine) newGoroutineWithRecover(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				var buf bytes.Buffer
				buf.Write(debug.Stack())
				log.Print(buf.String())
			}
		}()
		f()
	}()
}
