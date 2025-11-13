package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// simple implementation for periodic execution with semaphore.
// When the next schedule is run while the previous is still running, skip the tick
func tickerWithSemaphore() {
	ticker := time.NewTicker(2000 * time.Millisecond)

	sem := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	go func() {
		time.Sleep(20 * time.Second)
		log.Println("cancellation is triggered")
		cancel()
	}()

	for {
		select {
		case <-ticker.C:
			fmt.Println(">> ticker is ticking")
			// fmt.Println()
			select {
			// try to acquire the semaphore
			case sem <- struct{}{}:
				wg.Add(1)
				go doSomeLongWork(sem, &wg)
			default:
				// if previous work is still running when at the next 'tick'
				fmt.Println("previous job is still running, skip for this time")

			}
		case <-ctx.Done():
			log.Println("wait all the job to finish")
			wg.Wait()
			fmt.Println("done ...")
			return
		}
	}
}

func doSomeLongWork(sem <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-sem }()
	n := rand.Intn(3) + 1
	ticker := time.NewTicker(200 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Printf("wait for %d second\n", n)
		time.Sleep(time.Duration(n) * time.Second)
		cancel()
	}()
	for {
		select {
		case <-ticker.C:
			log.Println("||")
		case <-ctx.Done():
			ticker.Stop()
			log.Println("long work has done")
			fmt.Println("===================================")
			fmt.Println()
			return
		}
	}
}
