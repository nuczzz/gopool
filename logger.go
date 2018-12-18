package gopool

import "fmt"

type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type fmtLogger struct{}

func (f *fmtLogger) Error(args ...interface{}) {
	fmt.Println(args...)
}

func (f *fmtLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func newDefaultLogger() Logger {
	return &fmtLogger{}
}
