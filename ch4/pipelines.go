package main

import (
	"fmt"
	"math/rand"
)

// best practice for constructing pipelines

func pipelinesDemo() {
	generator := func(done <-chan any, integers ...int) <-chan int {
		intStream := make(chan int)

		go func() {
			defer close(intStream)
			for _, val := range integers {
				select {
				case <-done:
					fmt.Println("intStream done ...")
					return
				case intStream <- val:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan any, intStream <-chan int, multiplier int) <-chan int {
		multipierStream := make(chan int)
		go func() {
			defer close(multipierStream)
			for val := range intStream {
				select {
				case <-done:
					fmt.Println("multiply done ...")
					return
				case multipierStream <- val * multiplier:
				}
			}
		}()
		return multipierStream
	}

	add := func(done <-chan any, intStream <-chan int, additive int) <-chan int {
		addStream := make(chan int)
		go func() {
			defer close(addStream)
			for val := range intStream {
				select {
				case <-done:
					fmt.Println("add done ...")
					return
				case addStream <- val + additive:
				}
			}
		}()
		return addStream
	}

	done := make(chan any)
	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}

func generatorsDemo() {
	repeat := func(done <-chan any, values ...any) <-chan any {
		repeatStream := make(chan any)
		go func() {
			defer close(repeatStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						fmt.Println("repeatStream end ...")
						return
					case repeatStream <- v:
					}
				}
			}
		}()
		return repeatStream
	}

	take := func(done <-chan any, repeatStream <-chan any, num int) <-chan any {
		takeStream := make(chan any)
		go func() {
			defer close(takeStream)
			for range num {
				select {
				case <-done:
					fmt.Println("takeStream done ...")
					return
				case takeStream <- <-repeatStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan any)
	defer close(done)
	for num := range take(done, repeat(done, 1), 10) {
		fmt.Println(num)
	}
	fmt.Println("exit program ...")
}

func repeatFuncDemo() {
	take := func(done <-chan any, repeatStream <-chan any, num int) <-chan any {
		takeStream := make(chan any)
		go func() {
			defer close(takeStream)
			for range num {
				select {
				case <-done:
					fmt.Println("takeStream done ...")
					return
				case takeStream <- <-repeatStream:
				}
			}
		}()
		return takeStream
	}

	repeatFunc := func(done <-chan any, fn func() any) <-chan any {
		valueStream := make(chan any)
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					fmt.Println("repeatFunc end ...")
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	done := make(chan any)
	defer close(done)
	rand := func() any {
		return rand.Int()
	}
	for num := range take(done, repeatFunc(done, rand), 10) {
		fmt.Println(num)
	}
	fmt.Println("exit program")
}
