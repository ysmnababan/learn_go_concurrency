package main

import (
	"context"
	"log"
	"time"
)

func tickerWithCancel() {
	ticker := time.NewTicker(100 * time.Millisecond)
	now := time.Now()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	for {
		select {
		case <-ticker.C:
			log.Println("do some repetition: ", time.Since(now))

		case <-ctx.Done():
			ticker.Stop()
			log.Println("cancelled")
			cancel()
			return
		}
	}
}
