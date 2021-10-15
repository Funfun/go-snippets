package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	// Slice is a type that has reference to underlying array plus cap & len.
	a := []int{4567890, 456789} // a slice literal declaration, like a array
	b := make([]int, 1)
	copy(b, a)
	fmt.Println("a = ", a, len(a), cap(a))
	fmt.Println("b = ", b)

	// nil slice
	var nilSlice []int
	// appending to nil slice, creates new slice & assigns new slice to nilSlice
	nilSlice = append(nilSlice, a...) // apend one slice to another
	fmt.Println(nilSlice)

	// sort
	sort.Ints(a)
	fmt.Println(a)
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	fmt.Println(a)

	// args
	fmt.Println(os.Args)

}
