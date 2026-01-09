package main

import "fmt"

func main() {
	var a int = 10
	b := 2.1
	var c, max = 20, 30

	var d int64 = int64(c) // type conversions are explicit
	fmt.Printf("Hello,%d %f %10d %d Arch Linux!\n", a, b, d, max)
	var x, y = Sqrt(12.3)
	fmt.Printf("The answer of Sqrt 12.3 is: %.11f %.11f\n", x, y)

	primes := []int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	var s []int = primes[1:4]
	fmt.Println(s)

	OS()

	s1 := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
	}
	fmt.Println(s1)

	var a_map = map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
	}

	fmt.Println(a_map["B"])

	fmt.Println(WordCount("What the fuck is the that"))
}
