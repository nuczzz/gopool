package gopool

import (
	"bytes"
	"log"
	"runtime/debug"
)

type Task func()

type Goroutine interface {
	Execute(task Task)
	Terminal()
}

type goroutine struct {
	// pool who manage this goroutine.
	pool *goroutinePool

	// task channel for send or receive task, if task is nil,
	// release goroutine.
	task chan Task
}

// Execute execute task
func (g *goroutine) Execute(task Task) {
	g.task <- task
}

// Terminal terminal goroutine
func (g *goroutine) Terminal() {
	g.task <- nil
}

// run run goroutine
func (g *goroutine) run() {
	g.newGoroutineWithRecover(func() {
		for {
			select {
			case task := <-g.task:
				if task == nil {
					return
				}
				task()
				g.pool.recycleGoroutine(g)
			}
		}
	})
}

// newGoroutine create new goroutine and run it.
func newGoroutine(pool *goroutinePool) Goroutine {
	g := &goroutine{
		pool: pool,
		task: make(chan Task),
	}
	g.run()
	return g
}

// newGoroutineWithRecover create new goroutine with recover
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
