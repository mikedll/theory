package main

import (
	"fmt"
	"math/rand"
	"time"
)

func doWork(max int, resultsCh chan int, isDoneCh chan bool) {
	total := rand.Intn(max)
	
	for i := 1; i <= total; i++ {
		isDoneCh <- false
		resultsCh <- i
	}

	isDoneCh <- true
}

func main() {

	now := time.Now()
	rand.Seed(now.Unix())
	
	resultsCh := make(chan int)
	isDoneCh := make(chan bool)

	go doWork(100, resultsCh, isDoneCh)

	var isDone bool
	var result int
	isDone = <- isDoneCh
	for !isDone {
		result = <- resultsCh
		fmt.Printf("Result: %d\n", result)
		
		isDone = <-isDoneCh
	}
	
	fmt.Printf("Work is done\n")
}
