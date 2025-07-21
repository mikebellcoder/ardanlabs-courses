package main

import (
	"sync"
)

// counter is a variable incremented by all goroutines.
var counter int32

func main() {

	// Number of goroutines to use.
	const grs = 2

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(grs)

	var mu sync.Mutex

	// Create two goroutines.
	for range grs {
		go func() {
			for range 2 {

				// sync.Mutex is used to protect the counter variable.
				mu.Lock()
				{
					// Capture the value of Counter.
					value := counter
					// Increment our local value of Counter.
					value++
					// Store the value back into Counter.
					counter = value
				}
				mu.Unlock()

				/*
					// New version with atomic operations
					atomic.AddInt32(&counter, 1)
				*/

				/* Earlier version with data race
				// Capture the value of Counter.
				value := counter

				// Increment our local value of Counter.
				value++

				fmt.Println("logging")

				// Store the value back into Counter.
				counter = value
				*/

			}

			wg.Done()
		}()
	}

	wg.Wait()

	// Print the final value of Counter.
	println("Final Counter:", counter)
}
