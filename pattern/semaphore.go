package main

import (
	"fmt"
	"sync"
	"time"
)

// this is semaphore pattern for some application like db connection pooling
func semaphorePattern() {
	tokenCount := 5

	sem := make(chan struct{}, tokenCount)
	var wg sync.WaitGroup
	work := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()
		start := time.Now()

		sem <- struct{}{}
		fmt.Printf("Worker %d acquiring semaphore after waiting for %v \n", id, time.Since(start))

		time.Sleep(200 * time.Millisecond)
		fmt.Println("releasing semaphore", id)

		<-sem
	}

	for i := range 20 {
		wg.Add(1)
		go work(i, &wg)
	}
	wg.Wait()
	fmt.Println("done ....")
}
