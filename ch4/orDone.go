package main

import "fmt"

func orDoneDemo() {
	done := make(chan any)
	stream := make(chan any)
	go func() {
		for i := range 40 {
			stream <- i
		}
		close(stream)
	}()

	// instead of using this,
loop:
	for {
		select {
		case <-done:
			break loop
		case val, ok := <-stream:
			if !ok {
				return // or break
			}
			// do something with val
			_ = val
		}
	}

	// become simpler like this
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

	for val := range orDone(done, stream) {
		fmt.Println(val)
	}

	fmt.Println("demo done ...")
}
