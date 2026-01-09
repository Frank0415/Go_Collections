// split a string into words
package main

import (
	"strings"
)

func WordLen(s string) map[string]int {
	x := make(map[string]int)
	for _, substr := range strings.Fields(s) {
		x[substr] = len(substr)
	}
	return x
}

func WordCount(s string) map[string]int {
	x := make(map[string]int)
	for _, substr := range strings.Fields(s) {
		x[substr]++
	}
	return x
}
