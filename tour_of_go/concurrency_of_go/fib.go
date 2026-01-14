package main

import (
	"fmt"
)

func Fib_run() {
	value := make(chan int)
	quit := make(chan int)

	go func() {
		for range 10 {
			fmt.Println(<-value)
		}
		quit <- 0
	}()
	fib(value, quit)
}

func fib(value, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case value <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("From fib: return")
			return
		}
	}
}
