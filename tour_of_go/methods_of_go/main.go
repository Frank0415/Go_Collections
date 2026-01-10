package main

import (
	"fmt"
	"math"
)

type I interface {
	M()
}

type T struct {
	S string
}

func (t *T) String() string {
	return fmt.Sprintf("What is %s?", t.S)
}

func (t *T) M() {
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}

type F float64

func (f F) M() {
	fmt.Println(f)
}

// however there are cannot have methods to nil interfaces

func main() {
	var iii I

	var t *T
	iii = t
	describe(iii)
	iii.M()

	// tt := &T{"Hello"}
	iii = &T{"hello"}
	describe(iii)
	iii.M()

	fmt.Println("Using println,", iii)

	iii = F(math.Pi)
	describe(iii)
	iii.M()

	var i interface{} = "hello"

	s := i.(string)
	fmt.Println(s)

	s, ok := i.(string)
	fmt.Println(s, ok)

	f, ok := i.(float64)
	fmt.Println(f, ok)

	// f = i.(float64) // panic
	// fmt.Println(f)
}

func describe(i I) {
	fmt.Printf("(%v, %T)\n", i, i)
}
