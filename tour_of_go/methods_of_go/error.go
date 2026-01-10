package main

import (
	"fmt"
	"math"
	"time"
)

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s",
		e.When, e.What)
}

func Run() error {
	return &MyError{
		time.Now(),
		"it didn't work",
	}
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %.2f", e)
}

func Sqrt(x float64) (float64, error) {
	if x >= 0 {
		return math.Sqrt(x), nil
	} else {
		var err = ErrNegativeSqrt(x)
		return 0, err
	}
}
