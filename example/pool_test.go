package example

import (
	"fmt"
	"gopool"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := gopool.NewPoolWithDefault()

	var count int
	var lock sync.Locker

	expectMaxGoroutines := 1

	var wg sync.WaitGroup
	wg.Add(expectMaxGoroutines)
	for i := 0; i < expectMaxGoroutines; i++ {
		fmt.Printf("current goroutine: total[%v], working[%v], free[%v]\n",
			pool.GetTotalGoroutineNum(), pool.GetWorkingGoroutineNum(), pool.GetFreeGoroutineNum())
		pool.SubmitTask(func() {
			for j := 0; j < 10; j++ {
				lock.Lock()
				count++
				lock.Unlock()
			}
			wg.Done()
		})
	}
	wg.Wait()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	time.Sleep(5 * time.Second)
	fmt.Println("sum: ", count)
	fmt.Printf("current goroutine: total[%v], working[%v], free[%v]\n",
		pool.GetTotalGoroutineNum(), pool.GetWorkingGoroutineNum(), pool.GetFreeGoroutineNum())
}
