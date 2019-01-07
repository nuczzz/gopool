package gopool

type Pool interface {
	// GetTotalGoroutineNum return the sum of idle goroutine number
	// and working goroutine number.
	GetTotalGoroutineNum() int

	// GetFreeGoroutineNum return the number of idle goroutine.
	GetIdleGoroutineNum() int

	// GetWorkingGoroutineNum return the number of working goroutine.
	GetWorkingGoroutineNum() int

	// SubmitTask submit task to goroutine pool.
	SubmitTask(task func()) error
}

func NewPool(maxGoNum, maxIdleGoNum int) Pool {
	return newPool(maxGoNum, maxIdleGoNum)
}

func NewPoolWithDefault() Pool {
	return newPool(DefaultMaxGoroutineNum, DefaultMaxIdleGoroutineNum)
}
