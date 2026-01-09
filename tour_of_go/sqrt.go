package main

import (
	"fmt"
	"math"
	"runtime"
)

type Haha struct {
	a int32
	b bool
}

func Sqrt(x float64) (float64, float64) {
	const diff float64 = 1e-10
	var z float64 = x
	solution := math.Sqrt(x)
	i := 0
	for (z - solution) > diff {
		i++
		z -= (z*z - x) / (2 * z)
	}
	defer fmt.Println("After ", i, " iterations has completed.")
	// defer fmt.Println("")
	// defer fmt.Println("After ", i+1, " iterations has completed.")
	return z, solution
}

func OS() {
	fmt.Printf("The OS is ")
	switch os := runtime.GOOS; os {
	case "linux":
		fmt.Println("Linux.")
	case "darwin":
		fmt.Println("MACOS.")
	default:
		fmt.Printf("%s.", os)
	}
}

func myappend(vs... int)