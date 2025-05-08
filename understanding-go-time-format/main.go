package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	scn := bufio.NewScanner(os.Stdin)
	now := time.Now()

	for {
		fmt.Printf("> ")
		if !scn.Scan() {
			break
		}

		log.Println(now.Format(scn.Text()))
	}
}
