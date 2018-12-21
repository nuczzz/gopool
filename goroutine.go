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

// todo: when a goroutine free long time, kill it.
type goroutine struct {
	// pool who manage this goroutine
	pool *goroutinePool

	// task channel for send or receive task
	task chan Task
}

func (g *goroutine) Execute(task Task) {
	g.task <- task
}

func (g *goroutine) run() {
	newGoroutineWithRecover(func() {
		for f := range g.task {
			f()
			g.pool.recycleGoroutine(g)
		}
	})
}

func newGoroutine(pool *goroutinePool) Goroutine {
	g := &goroutine{
		pool: pool,
		task: make(chan Task, 1),
	}
	g.run()
	return g
}

func newGoroutineWithRecover(f func()) {
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
