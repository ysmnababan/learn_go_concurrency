package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, jobs <-chan int) {
	for {
		select {
		case i := <-jobs:
			fmt.Println(i)
		case <-ctx.Done():
			fmt.Println("cancelled")
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func simpleSelectPattern() {
	ctx, cancel := context.WithCancel(context.Background())

	jobs := make(chan int)
	go worker(ctx, jobs)
	go func() {
		for i := range 8 {
			jobs <- i
		}
		close(jobs)
	}()

	fmt.Println("before sleep")

	time.Sleep(500 * time.Millisecond)
	cancel()
	time.Sleep(500 * time.Millisecond)
	fmt.Println("end")
}
