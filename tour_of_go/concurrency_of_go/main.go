package main

import (
	"fmt"
	"sync"
)

func runBackground(wg *sync.WaitGroup, fn func()) {
    wg.Go(func() {
        fn()
    })
}

func say(s string) {
	ctr := 0
	for range 10 {
		// time.Sleep(time.Millisecond)
		// fmt.Println(s)
		fmt.Println(s, " ", ctr)
		ctr++
	}

}

func main() {
	var wg sync.WaitGroup
	
    runBackground(&wg, func() { say("goroutine 1") })
    runBackground(&wg, func() { say("goroutine 2") })
    runBackground(&wg, func() { Fib_run() })

	wg.Wait()
}
