package main

import (
	"fmt"
	"time"
)

func main() {
	go process()
	fmt.Println("next")
	time.Sleep(2 * time.Second) // wait for goroutine finish
}

func process() {
	time.Sleep(1 * time.Second) // heavy process
	fmt.Println("goroutine completed")
}
