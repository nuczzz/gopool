package gopool

import (
	"sync"
)

type goroutinePool struct {
	// maxGoroutineNum max number of goroutine.the number of goroutines
	// in pool is freeGoroutineNum+workingGoroutineNum.
	maxGoroutineNum int

	// lock lock of goroutine pool.
	lock sync.Mutex

	// freeGoroutines slice of free goroutines.if freeGoroutineNum is
	// GREATER THAN ZERO, we get the last free Goroutine in freeGoroutines
	// to handle a new request,and remove the Goroutine from freeGoroutines.
	// if freeGoroutineNum is ZERO and total goroutines LESS THAN maxGoroutineNum,
	// we create a new Goroutine to handle it.after a Goroutine handled a request,
	// put the Goroutine into freeGoroutines.
	freeGoroutines []Goroutine

	// freeGoroutineNum number of free goroutines in freeGoroutines.
	freeGoroutineNum int32

	// workingGoroutineNum number of working goroutine.
	workingGoroutineNum int32

	// logger of goroutine pool.
	logger Logger
}

func (gp *goroutinePool) getGoroutine() (Goroutine, error) {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	if gp.freeGoroutineNum > 0 {
		g := gp.freeGoroutines[gp.freeGoroutineNum-1]
		gp.freeGoroutines = gp.freeGoroutines[:gp.freeGoroutineNum-1]
		gp.freeGoroutineNum--
		gp.workingGoroutineNum++
		return g, nil
	} else if !gp.isOverflow() {
		g := newGoroutine(gp)
		gp.workingGoroutineNum++
		return g, nil
	}

	return nil, ErrPoolOverflow
}

func (gp *goroutinePool) isOverflow() bool {
	return gp.freeGoroutineNum+gp.workingGoroutineNum >= int32(gp.maxGoroutineNum)
}

func (gp *goroutinePool) recycleGoroutine(g Goroutine) {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	gp.workingGoroutineNum--
	gp.freeGoroutineNum++
	gp.freeGoroutines = append(gp.freeGoroutines, g)
}

func (gp *goroutinePool) GetTotalGoroutineNum() int {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	return int(gp.freeGoroutineNum + gp.workingGoroutineNum)
}

func (gp *goroutinePool) GetFreeGoroutineNum() int {
	return int(gp.freeGoroutineNum)
}

func (gp *goroutinePool) GetWorkingGoroutineNum() int {
	return int(gp.workingGoroutineNum)
}

func (gp *goroutinePool) SubmitTask(task func()) error {
	goroutine, err := gp.getGoroutine()
	if err != nil {
		return nil
	}
	goroutine.Execute(task)
	return nil
}

func newPool(max int, logger Logger) Pool {
	return &goroutinePool{
		maxGoroutineNum: max,
		logger:          logger,
		freeGoroutines:  make([]Goroutine, 0, max),
	}
}
