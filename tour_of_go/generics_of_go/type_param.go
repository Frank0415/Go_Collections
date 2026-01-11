package main

import "fmt"

func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		// v and x are type T, which has the comparable
		// constraint, so we can use == here.
		if v == x {
			return i
		}
	}
	return -1
}

func Demo_index() {
	// Index works on a slice of ints
	si := []int{10, 20, 15, -10}
	fmt.Println(Index(si, 15))

	// Index also works on a slice of strings
	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(ss, "hello"))
}

type linked_list[T any] struct {
	next *linked_list[T]
	val  T
}

func Demo_Linked_List() {
	head := &linked_list[int]{val: 10}
	current := head
	for i := 0; i < 5; i++ {
		current.next = &linked_list[int]{val: current.val + 10}
		current = current.next
	}

	for n := head; n != nil; n = n.next {
		fmt.Printf("%v ", n.val)
	}
	fmt.Println()
}
