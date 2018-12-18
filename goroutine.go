package gopool

type Goroutine interface {
	Execute(task func())
}

type goroutine struct {
	// pool who manage this goroutine
	pool *goroutinePool

	// task channel for send or receive task
	task chan func()
}

func (g *goroutine) Execute(task func()) {
	g.task <- task
}

func (g *goroutine) run() {
	g.newGoroutineWithRecover(func() {
		for f := range g.task {
			f()
		}
	})
}

func newGoroutine() Goroutine {
	g := &goroutine{
		task: make(chan func(), 1),
	}
	g.run()
	return g
}

func (g *goroutine) newGoroutineWithRecover(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				g.pool.logger.Error(err)
				return
			}
		}()
		f()
		g.pool.recycleGoroutine(g)
	}()
}
