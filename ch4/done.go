package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Prevent goroutine channel from leaking using 'done' channel.
// Also you can see that the 'terminated' chan can be blocking, and by closing it, the program will continue
func demoDone() {
	doWork := func(done <-chan any, strings <-chan string) <-chan any {
		terminated := make(chan any)
		go func() {
			defer fmt.Println("work is done.")
			defer close(terminated)

			for {
				select {
				case <-done:
					return
				case s := <-strings:
					fmt.Println(s)
				}
			}
		}()
		return terminated
	}

	done := make(chan any)
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("closing the 'do work' func")
		close(done)
	}()
	fmt.Println("before the <- terminated")
	<-terminated
	fmt.Println("Done.")
}

// We can also create a channel for streaming random values.
// We can achieve this by sending it to channel and read it until the signal 'done' is sent.
func randStreamDemo() {
	newRandStream := func(done <-chan any) <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("stream has completed.")
			defer close(randStream)
			for {
				select {
				case <-done:
					return
				default:
					randStream <- rand.Int()
				}
			}
		}()

		return randStream
	}

	done := make(chan any)

	randStream := newRandStream(done)

	i := 1
	for val := range randStream {
		fmt.Printf("%d: %d\n", i, val)
		if i == 3 {
			fmt.Println("close the channel")
			close(done)
		}
		i++
	}
	fmt.Println("exit program ... ")
}
