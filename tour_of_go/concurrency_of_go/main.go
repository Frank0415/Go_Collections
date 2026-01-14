package main

import (
	"fmt"
	"time"
)

func say(s string) {
	ctr := 0
	for range 500 {
		time.Sleep(time.Millisecond)
		// fmt.Println(s)
		fmt.Println(s, " ", ctr)
		ctr++
	}

}

func main() {
	go say("goroutine 1")
	go say("goroutine 2")
	say("goroutine main")
}
