package main

import (
	"crypto/sha1"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

func init() {
	// Allocate one logical processor for scheduler to use.
	runtime.GOMAXPROCS(1)
}

func main() {

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Create Goroutines")

	// Create the first goroutine and manage its lifecycle here.
	go func() {
		printHashes("A")
		wg.Done()
	}()

	// Create the second goroutine and manage its lifecycle here.
	go func() {
		printHashes("B")
		wg.Done()
	}()

	// Wait for the goroutines to finish.
	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")
}

// printHashes prints a sequence of hashes with a given prefix.
func printHashes(prefix string) {
	// print each hash from 1 to 10. Change this to 5000 and
	// see how the scheduler behaves.
	for i := 1; i <= 5000; i++ {
		// Convert i to a string.
		num := strconv.Itoa(i)
		// Calculate the hash for string num.
		sum := sha1.Sum([]byte(num))
		// Print prefix: 5-digit-number: hex encoded hash.
		fmt.Printf("%s: %05d: %x\n", prefix, i, sum)
	}

	fmt.Println("Completed")
}
