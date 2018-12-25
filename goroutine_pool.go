package gopool

import (
	"errors"
	"sync"
	"time"
)

const (
	DefaultMaxGoroutineNum = 50000
	//DefaultMaxLifeTime     = 5 * 60 //second
	DefaultMaxLifeTime     = 3 //second
)

var ErrPoolOverflow = errors.New("pool overflow")

type goroutinePool struct {
	// maxGoroutineNum max number of goroutine.the number of goroutines
	// in pool is idleGoroutineNum+workingGoroutineNum.
	maxGoroutineNum int

	// lock lock of goroutine pool.
	lock sync.Mutex

	// idleGoroutines slice of idle goroutines.if idleGoroutineNum is
	// GREATER THAN ZERO, we get the last idle Goroutine in idleGoroutines
	// to handle a new request,then remove the Goroutine from idleGoroutines.
	// if idleGoroutineNum is ZERO and total goroutines LESS THAN maxGoroutineNum,
	// we create a new Goroutine to handle it.after a Goroutine handled a request,
	// put the Goroutine into idleGoroutines.
	idleGoroutines []Goroutine

	// idleGoroutineNum number of idle goroutines in idleGoroutines.
	idleGoroutineNum int32

	// workingGoroutineNum number of working goroutine.
	workingGoroutineNum int32

	// maxLifeTime max life time of idle goroutine.
	maxLifeTime time.Duration
}

func (gp *goroutinePool) getGoroutine() (Goroutine, error) {
	gp.lock.Lock()
	defer gp.lock.Unlock()
	if gp.idleGoroutineNum > 0 {
		g := gp.idleGoroutines[gp.idleGoroutineNum-1]
		gp.idleGoroutines = gp.idleGoroutines[:gp.idleGoroutineNum-1]
		gp.idleGoroutineNum--
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
	return gp.idleGoroutineNum+gp.workingGoroutineNum >= int32(gp.maxGoroutineNum)
}

func (gp *goroutinePool) recycleGoroutine(g Goroutine) {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	gp.workingGoroutineNum--
	gp.idleGoroutineNum++
	gp.idleGoroutines = append(gp.idleGoroutines, g)
	g.ResetTimeout()
}

// todo: goroutine pool check idle goroutine timeout periodically maybe better
func (gp *goroutinePool) releaseGoroutine(g Goroutine) {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	for i := range gp.idleGoroutines {
		if g == gp.idleGoroutines[i] {
			gp.idleGoroutines = append(gp.idleGoroutines[:i], gp.idleGoroutines[i+1:]...)
			gp.idleGoroutineNum--
			return
		}
	}
}

func (gp *goroutinePool) GetTotalGoroutineNum() int {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	return int(gp.idleGoroutineNum + gp.workingGoroutineNum)
}

func (gp *goroutinePool) GetIdleGoroutineNum() int {
	return int(gp.idleGoroutineNum)
}

func (gp *goroutinePool) GetWorkingGoroutineNum() int {
	return int(gp.workingGoroutineNum)
}

// SubmitTask submit task to goroutine pool, goroutine pool
// will try to get a free goroutine or create a new goroutine.
// Returns a error if the count of goroutine pool current
// goroutines over maxGoroutineNum. If get a goroutine but
// when the goroutine execute task return error(stop goroutine
// timer failed), we will try to get another goroutine to
// execute the task again.
func (gp *goroutinePool) SubmitTask(task func()) error {
	for {
		goroutine, err := gp.getGoroutine()
		if err != nil {
			return err
		}
		if goroutine.Execute(task) == nil {
			return nil
		}
	}
}

func newPool(max int) Pool {
	return &goroutinePool{
		maxGoroutineNum: max,
		idleGoroutines:  make([]Goroutine, 0, max),
		maxLifeTime:     time.Second * time.Duration(DefaultMaxLifeTime),
	}
}
