package main

import "fmt"

func fibonacci() func() int {
	var val, prev_val int
	val = 1
	return func() int {
		var tmp = val
		var tmp2 = prev_val
		val += prev_val
		prev_val = tmp
		return tmp2
	}
}

func Fibonacci_test() {
	var fib = fibonacci()
	for range 10 {
		fmt.Println(fib())
	}
}
