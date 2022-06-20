package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func worker(ctx context.Context) {
LABEL:
	for {
		select {
		case <-ctx.Done():
			break LABEL
		default:
			fmt.Println("worker...")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go worker(ctx)
	cancel()
	time.Sleep(time.Second * 2)
	fmt.Println("over")
}
