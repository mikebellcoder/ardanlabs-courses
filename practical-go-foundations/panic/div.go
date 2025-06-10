package main

import "fmt"

func main() {
	safeDiv(4, 0)
}

func safeDiv(a, b int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in safeDiv", r)
		}
	}()

	return a / b, nil
}
