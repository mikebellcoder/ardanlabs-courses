package main

import (
	"fmt"
	"time"
)

func main() {
	go fmt.Println("goroutine")
	fmt.Println("main")
	for i := range 3 {
		go func() {
			fmt.Println("goroutine", i)
		}()
	}

	time.Sleep(10 * time.Millisecond)
}
