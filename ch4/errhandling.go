package main

import (
	"fmt"
	"log"
	"net/http"
)

// in go, error should be considered first class citizens
// so to handle error, you can send it along with the result
func errorHandlingDemo() {
	type Result struct {
		Response *http.Response
		Error    error
	}
	checkStatus := func(done <-chan any, urls []string) <-chan Result {
		res := make(chan Result)
		go func() {
			defer close(res)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result.Error = err
				result.Response = resp
				select {
				case <-done:
					fmt.Println("checkStatus done")
					return
				case res <- result:
				}
			}
		}()
		return res
	}

	done := make(chan any)

	urls := []string{"https://www.google.com", "https://badhost"}
	for resp := range checkStatus(done, urls) {
		if resp.Error != nil {
			log.Println(resp.Error)
			continue
		}
		fmt.Println(resp.Response.Status)
	}
}
