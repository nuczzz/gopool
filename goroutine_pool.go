package gopool

import (
	"errors"
	"sync"
	"time"
)

const (
	DefaultMaxGoroutineNum     = 50000
	DefaultMaxIdleGoroutineNum = 10
	//DefaultCleanPeriod     = 5 * 60 //second
	DefaultCleanPeriod = 3 //second
)

var ErrPoolOverflow = errors.New("pool overflow")

type goroutinePool struct {
	// maxGoroutineNum max number of goroutine.the number of goroutines
	// in pool is idleGoroutineNum+workingGoroutineNum.
	maxGoroutineNum int

	// maxIdleGoroutineNum max number of idle goroutine.
	maxIdleGoroutineNum int

	// lock lock of goroutine pool.
	lock sync.Mutex

	// idleGoroutines slice of idle goroutines.if idleGoroutineNum is
	// GREATER THAN ZERO, we get the last idle Goroutine in idleGoroutines
	// to handle a new request,then remove the Goroutine from idleGoroutines.
	// if idleGoroutineNum is ZERO and total goroutines LESS THAN maxGoroutineNum,
	// we create a new Goroutine to handle it.after a Goroutine handled a request,
	// put the Goroutine into idleGoroutines.
	idleGoroutines []Goroutine

	// workingGoroutineNum number of working goroutine.
	workingGoroutineNum int

	// ticker ticker to release idle goroutine.
	ticker *time.Ticker
}

func (gp *goroutinePool) getGoroutine() (Goroutine, error) {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	length := len(gp.idleGoroutines)
	if length > 0 {
		g := gp.idleGoroutines[length-1]
		gp.idleGoroutines[length-1] = nil
		gp.idleGoroutines = gp.idleGoroutines[:length-1]
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
	return len(gp.idleGoroutines)+gp.workingGoroutineNum >= gp.maxGoroutineNum
}

func (gp *goroutinePool) recycleGoroutine(g Goroutine) {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	gp.workingGoroutineNum--
	gp.idleGoroutines = append(gp.idleGoroutines, g)
}

// cleanGoroutinePeriodically clean idle goroutine if idleGoroutineNum greater
// than maxIdleGoroutineNum.
func (gp *goroutinePool) cleanGoroutinePeriodically() {
	for {
		select {
		case <-gp.ticker.C:
			gp.lock.Lock()
			length := len(gp.idleGoroutines)
			if length > gp.maxIdleGoroutineNum {
				for i := gp.maxIdleGoroutineNum; i < length; i++ {
					gp.idleGoroutines[i].Terminal()
					gp.idleGoroutines[i] = nil
				}
				gp.idleGoroutines = gp.idleGoroutines[:gp.maxIdleGoroutineNum]
			}
			gp.lock.Unlock()
		}
	}
}

func (gp *goroutinePool) GetTotalGoroutineNum() int {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	return len(gp.idleGoroutines) + gp.workingGoroutineNum
}

func (gp *goroutinePool) GetIdleGoroutineNum() int {
	return len(gp.idleGoroutines)
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
	goroutine, err := gp.getGoroutine()
	if err != nil {
		return err
	}
	goroutine.Execute(task)
	return nil
}

func newPool(maxGoNum, maxIdleGoNum int) Pool {
	if maxGoNum <= 0 {
		maxGoNum = DefaultMaxGoroutineNum
	}
	if maxIdleGoNum <= 0 {
		maxIdleGoNum = DefaultMaxIdleGoroutineNum
	}
	pool := &goroutinePool{
		maxGoroutineNum:     maxGoNum,
		maxIdleGoroutineNum: maxIdleGoNum,
		idleGoroutines:      make([]Goroutine, 0, maxGoNum),
		ticker:              time.NewTicker(time.Second * DefaultCleanPeriod),
	}
	go pool.cleanGoroutinePeriodically()
	return pool
}
