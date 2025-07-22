package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"
)

func main() {
	// waitForResult()
	// fanOut()
	// waitForTask()
	// pooling()
	// fanOutSem()
	// fanOutBounded()
	// selectExample()
	// drop()
	cancellation()
}

func waitForResult() {
	ch := make(chan string)

	go func() {
		time.Sleep(time.Duration(rand.IntN(500)) * time.Millisecond)
		ch <- "paper"
		fmt.Println("employee : send signal")
	}()

	p := <-ch
	fmt.Println("manager : received signal:", p)

	time.Sleep(time.Second)
	fmt.Println("------------------------------------------------------------------")
}

func fanOut() {
	emps := 2_000
	ch := make(chan string, emps)

	for e := range emps {
		go func(emp int) {
			time.Sleep(time.Duration(rand.IntN(200)) * time.Millisecond)
			ch <- "paper"
			fmt.Println("employee : sent signal :", emp)
		}(e)
	}

	for emps > 0 {
		p := <-ch
		emps--
		fmt.Println(p)
		fmt.Println("manager : received signal :", emps)
	}

	time.Sleep(time.Second)
	fmt.Println("------------------------------------------------------------------")
}

func waitForTask() {
	ch := make(chan string)

	go func() {
		p := <-ch
		fmt.Println("employee : received signal :", p)
	}()

	time.Sleep(time.Duration(rand.IntN(500)) * time.Millisecond)
	ch <- "paper"
	fmt.Println("manager : sent signal")

	time.Sleep(time.Second)
	fmt.Println("------------------------------------------------------------------")
}

func pooling() {
	ch := make(chan string)

	g := runtime.NumCPU()
	for e := range g {
		go func(emp int) {
			for p := range ch {
				fmt.Printf("employee %d : received signal: %s\n", emp, p)
			}
			fmt.Printf("employee %d : received shutdown signal\n", emp)
		}(e)
	}

	const work = 100
	for w := range work {
		ch <- "paper"
		fmt.Println("manager : sent signal :", w)
	}

	close(ch)
	fmt.Println("manager : sent shutdown signal")
}

func fanOutSem() {
	emps := 2_000
	ch := make(chan string, emps)

	g := runtime.NumCPU()
	sem := make(chan bool, g)

	for e := range emps {
		go func(emp int) {
			sem <- true
			{
				time.Sleep(time.Duration(rand.IntN(200)) * time.Millisecond)
				ch <- "paper"
				fmt.Println("employee : sent signal :", emp)
			}
			<-sem
		}(e)
	}

	for emps > 0 {
		p := <-ch
		emps--
		fmt.Println(p)
		fmt.Println("manager : received signal :", emps)
	}

	time.Sleep(time.Second)
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("Number of NumCPU/Goroutines runnable state:", g)
}

func fanOutBounded() {
	work := []string{"paper", "paper", "paper", "paper", "paper", 2000: "paper"}

	g := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for e := range g {
		go func(emp int) {
			defer wg.Done()
			for p := range ch {
				fmt.Printf("employee %d : recv'd signal : %s\n", emp, p)
			}
			fmt.Printf("employee %d : recv'd shutdown signal\n", emp)
		}(e)
	}

	for _, wrk := range work {
		ch <- wrk
	}
	close(ch)
	wg.Wait()

	fmt.Println("----------------------------------------------------------------")

}

func selectExample() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "one"
	}()

	for range 2 {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
			// default:
			// fmt.Println("defaulted!")
		}
	}
}

func drop() {
	const cap = 100
	ch := make(chan string, cap)

	go func() {
		for p := range ch {
			fmt.Println("employee : recv'd signal :", p)
		}
	}()

	const work = 2000
	for w := range work {
		select {
		case ch <- "paper":
			fmt.Println("manager : sent signal :", w)
		default:
			fmt.Println("manager : dropped data :", w)
		}
	}

	close(ch)
	fmt.Println("manager : sent shutdown signal")

	time.Sleep(time.Second)
	fmt.Println("-----------------------------------------------------------------------------")
}

func cancellation() {
	duration := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	ch := make(chan string, 1)

	go func() {
		time.Sleep(time.Duration(rand.IntN(200)) * time.Millisecond)
		ch <- "paper"
	}()

	select {
	case d := <-ch:
		fmt.Println("work complete", d)

	case <-ctx.Done():
		fmt.Println("work completed")
	}

	time.Sleep(time.Second)
	fmt.Println("-----------------------------------------------------------------------------")
}
