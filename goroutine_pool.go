package gopool

type GoroutinePool struct {
	//coreGoroutineNum is the number of core goroutine.
	//we create a new goroutine to handle the new request
	//when goroutine number less than coreGoroutineNum,
	//even if some goroutines are not working.
	coreGoroutineNum int

	//maxGoroutineNum is the max number of goroutine.
	//maxGoroutineNum count all of the goroutines(include core goroutine).
	maxGoroutineNum int
}
