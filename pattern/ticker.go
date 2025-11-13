package main

import (
	"fmt"
	"time"
)

func simpleTickerPattern() {
	ticker := time.NewTicker(500 * time.Millisecond)

	defer ticker.Stop()

	// you can use this pattern if there is only 1 channel
	for range ticker.C {
		fmt.Println(time.Now().UTC())
	}

}
