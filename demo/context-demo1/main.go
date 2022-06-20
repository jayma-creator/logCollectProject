package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var exit bool

func worker() {
	defer wg.Done()
	for {
		fmt.Println("worker...")
		time.Sleep(time.Second)
		if exit {
			break
		}
	}
}

func main() {
	wg.Add(1)
	go worker()
	time.Sleep(time.Second * 5)
	exit = true
	wg.Wait()
	fmt.Println("over")
}
