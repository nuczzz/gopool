package gopool

import (
	"errors"
)

const (
	DefaultMaxGoroutineNum = 50000
)

var (
	ErrPoolOverflow = errors.New("pool overflow")

	defaultLogger = newDefaultLogger()
)
