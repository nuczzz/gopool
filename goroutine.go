package gopool

import (
	"bytes"
	"errors"
	"log"
	"runtime/debug"
	"time"
)

type Task func()

type Goroutine interface {
	Execute(task Task) error
	ResetTimeout()
}

type goroutine struct {
	// pool who manage this goroutine.
	pool *goroutinePool

	// task channel for send or receive task.
	task chan Task

	// timer timer for release idle goroutine resource.
	timer *time.Timer
}

// Execute stop the timer before a goroutine execute task,
// and returns error if stop timer failed(timer has stopped).
// When g.timer received timeout signal, goroutine pool will
// release the goroutine, and the task will not be executed
// forever.
func (g *goroutine) Execute(task Task) error {
	if !g.timer.Stop() {
		return errors.New("stop timer false")
	}
	g.task <- task
	return nil
}

func (g *goroutine) ResetTimeout() {
	g.timer.Reset(g.pool.maxLifeTime)
}

func (g *goroutine) run() {
	g.newGoroutineWithRecover(func() {
		for {
			select {
			case task := <-g.task:
				task()
				g.pool.recycleGoroutine(g)
			case <-g.timer.C:
				goto forEnd
			}
		}
	forEnd:
		g.pool.releaseGoroutine(g)
	})
}

func newGoroutine(pool *goroutinePool) Goroutine {
	g := &goroutine{
		pool:  pool,
		task:  make(chan Task),
		timer: time.NewTimer(pool.maxLifeTime),
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
