package gopool

type Pool interface {
	//Get a goroutine from goroutine pool
	Get() (Goroutine, error)
}

func NewPool(maxNum, coreNum int) Pool {
	if coreNum > maxNum {
		coreNum = maxNum
	}

	return &goroutinePool{
		maxGoroutineNum:  maxNum,
		coreGoroutineNum: coreNum,
		goroutines:       make([]Goroutine, 0, maxNum),
	}
}
