package it_test

import (
	"fmt"
	"strconv"

	"github.com/gomoni/it"
)

func ExampleChain_Filter() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	res := it.NewChain(it.From(n)).
		Filter(func(s string) bool { return len(s) >= 2 }).
		Filter(func(s string) bool { return len(s) >= 3 }).
		Slice()
	fmt.Println(res)
	// Output: [aaa aaaaaaa]
}

func ExampleChain_Reduce() {
	m := []int{1, 2, 3, 4, 5, 6, 7}
	count := it.NewChain(it.From(m)).
		Reduce(func(a, _ int) int { return a + 1 }, 0)
	fmt.Println(count)
	// Output: 7
}

func ExampleMappable_Map() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	// maps string->int
	res := it.NewMappable[string, int](it.From(n)).
		Map(func(s string) int { return len(s) }).
		Filter(func(i int) bool { return i >= 2 }).
		Slice()
	fmt.Println(res)
	// Output: [2 3 7]
}

func ExampleMappable_Map_second() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	// maps string->int and int->string
	res := it.NewMappable[string, int](it.From(n)).
		Map(func(s string) int { return len(s) }).
		Filter(func(i int) bool { return i >= 2 }).
		Map(func(i int) string { return "string(" + strconv.Itoa(i) + ")" }).
		Slice()
	fmt.Println(res)
	// Output: [string(2) string(3) string(7)]
}
