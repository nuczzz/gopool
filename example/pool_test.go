package example

import (
	"fmt"
	"github.com/nuczzz/gopool"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	pool := gopool.NewPoolWithDefault()

	var count int
	var lock sync.Mutex

	expectMaxGoroutines := 100

	var wg sync.WaitGroup
	wg.Add(expectMaxGoroutines)
	for i := 0; i < expectMaxGoroutines; i++ {
		fmt.Printf("current goroutine: total[%v], working[%v], free[%v]\n",
			pool.GetTotalGoroutineNum(), pool.GetWorkingGoroutineNum(), pool.GetFreeGoroutineNum())
		task := func() {
			for j := 0; j < 100; j++ {
				lock.Lock()
				count++
				lock.Unlock()
			}
			wg.Done()
		}
		if err := pool.SubmitTask(task); err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait()
	fmt.Println("sum: ", count)
	fmt.Printf("current goroutine: total[%v], working[%v], free[%v]\n",
		pool.GetTotalGoroutineNum(), pool.GetWorkingGoroutineNum(), pool.GetFreeGoroutineNum())
}
