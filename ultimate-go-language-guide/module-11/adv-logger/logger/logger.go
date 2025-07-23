package logger

import (
	"fmt"
	"io"
	"sync"
)

// Logger is a custom logger with ability to use drop pattern when issues arrise from io.
type Logger struct {
	ch chan string
	wg sync.WaitGroup
}

// New constructor for Logger
func New(w io.Writer, cap int) *Logger {
	l := Logger{
		ch: make(chan string, cap),
	}

	l.wg.Add(1)

	go func() {
		l.wg.Done()
		for v := range l.ch {
			fmt.Fprintln(w, v)
		}
	}()

	return &l
}

// Shutdown ...
func (l *Logger) Shutdown() {
	close(l.ch)
	l.wg.Wait()
}

// Println ...
func (l *Logger) Println(v string) {
	select {
	case l.ch <- v:
	default:
		fmt.Println("DROP")
	}
}
