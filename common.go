package gopool

import (
	"fmt"
)

func NewGoroutineWithRecover(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				return
			}
		}()

		f()
	}()
}
