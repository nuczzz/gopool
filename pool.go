package gopool

type Pool interface {
	// GetTotalGoroutineNum return the sum of free goroutine number
	// and working goroutine number.
	GetTotalGoroutineNum() int

	// GetFreeGoroutineNum return the number of free goroutine.
	GetFreeGoroutineNum() int

	// GetWorkingGoroutineNum return the number of working goroutine.
	GetWorkingGoroutineNum() int

	// SubmitTask submit task to goroutine pool.
	SubmitTask(task func()) error
}

func NewPool(max int) Pool {
	return newPool(max)
}

func NewPoolWithDefault() Pool {
	return newPool(DefaultMaxGoroutineNum)
}
