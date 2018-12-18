package gopool

type Pool interface {
	//Get a goroutine from goroutine pool
	GetGoroutine() (Goroutine, error)

	GetWorkingGoroutineNum() int

	GetTotalGoroutineNum() int
}

func NewPool(max int) Pool {
	return newPool(max, defaultLogger)
}

func NewPoolWithDefault() Pool {
	return newPool(DefaultMaxGoroutineNum, defaultLogger)
}

func NewPoolWithLogger(max int, logger Logger) Pool {
	return newPool(max, logger)
}
