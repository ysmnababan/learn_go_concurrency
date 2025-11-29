package main

import (
	"fmt"
	"sync"
)

func teeDemo() {
	done := make(chan any)
	defer close(done)

	orDone := func(done <-chan any, streamVal <-chan any) <-chan any {
		stream := make(chan any)

		go func() {
			defer close(stream)
			for {
				select {
				case <-done:
					return
				case val, ok := <-streamVal:
					if !ok {
						return
					}
					select {
					case stream <- val:
					case <-done:
					}
				}
			}
		}()
		return stream
	}
	tee := func(done <-chan any, in <-chan any) (_, _ <-chan any) {
		out1, out2 := make(chan any), make(chan any)
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
				o1, o2 := out1, out2 // copy first so it wont affect the original channel
				for range 2 {
					select {
					case <-done:
					case o1 <- val:
						o1 = nil // make it nil so it won't be selected next time,
					case o2 <- val:
						o2 = nil
					}
				}
			}
		}()
		return out1, out2
	}
	streamVal := make(chan any)
	go func() {
		for i := range 30 {
			streamVal <- i
		}
		close(streamVal)
	}()
	out1, out2 := tee(done, streamVal)

	var wg sync.WaitGroup
	wg.Add(2)

	process := func(in <-chan any, name string, wg *sync.WaitGroup) {
		defer wg.Done()
		for val := range in {
			fmt.Printf("from %s: %d\n", name, val)
		}
	}

	go process(out1, "out1", &wg)
	go process(out2, "out2", &wg)
	wg.Wait()
	fmt.Println("demo done ... ")
}
