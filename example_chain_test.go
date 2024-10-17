package it_test

import (
	"fmt"
	"slices"

	"github.com/gomoni/it"
)

func ExampleChain() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	ch := it.NewChain(slices.Values(n))
	slice := ch.
		Filter(func(s string) bool { return len(s) >= 2 }).
		Filter(func(s string) bool { return len(s) <= 4 }).
		Collect()
	fmt.Println(slice)
	// Output: [aa aaa]
}

func ExampleMappable() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	ch := it.NewMappable[string, int](slices.Values(n))
	slice := ch.
		Filter(func(s string) bool { return len(s) >= 2 }).
		Map(func(s string) int { return len(s) }).
		Collect()
	fmt.Println(slice)
	// Output: [2 3 7]
}
