package main

import (
	"fmt"
	"runtime"
	"sync"
)

func init() {

	// Allocate one logical processor for the program.
	runtime.GOMAXPROCS(1)
}

func main() {

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Start Goroutines")

	// Create a goroutine from the lowercase function.
	go func() {
		lowercase()
		fmt.Println("Lowercase goroutine finished")
		wg.Done()
	}()

	// Create a goroutine from the uppercase function.
	go func() {
		uppercase()
		fmt.Println("Uppercase goroutine finished")
		wg.Done()
	}()

	// Wait for the goroutines to finish.
	fmt.Println("Waiting To Finish")
	wg.Wait()
	// runtime.Gosched() // Yield the processor to allow other goroutines to run

	fmt.Println("\nTerminating Program")

}

// lowercase prints the lowercase alphabet.
func lowercase() {
	for range 2 {
		for c := 'a'; c <= 'z'; c++ {
			fmt.Printf("%c ", c)
		}
	}
	fmt.Println() // Print a newline after lowercase letters
}

// uppercase prints the uppercase alphabet.
func uppercase() {
	for range 2 {
		for c := 'A'; c <= 'Z'; c++ {
			fmt.Printf("%c ", c)
		}
	}
	fmt.Println() // Print a newline after uppercase letters
}
