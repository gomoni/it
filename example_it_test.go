package it_test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gomoni/it"
)

func ExampleFilter() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	s0 := it.From(n)
	s1 := it.Filter(s0, func(s string) bool { return len(s) >= 2 })
	slice := it.Slice(s1)
	fmt.Println(slice)
	// Output: [aa aaa aaaaaaa]
}

func ExampleMap() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	// maps string->int->float32
	s0 := it.From(n)
	s1 := it.Map(s0, func(s string) int { return len(s) })
	s2 := it.Map(s1, func(i int) float32 { return float32(i) })
	s3 := it.Map(s2, func(f float32) string { return strconv.FormatFloat(float64(f), 'E', 4, 32) })
	slice := it.Slice(s3)
	fmt.Println(slice)
	// Output: [2.0000E+00 3.0000E+00 7.0000E+00 1.0000E+00]
}

func ExampleSorted() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	s0 := it.From(n)
	s1 := it.Sorted(s0, strings.Compare)
	slice := it.Slice(s1)
	fmt.Println(slice)
	// Output: [a aa aaa aaaaaaa]
}

func ExampleReverse() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	s0 := it.From(n)
	s1 := it.Sorted(s0, strings.Compare)
	s2 := it.Reverse(s1)
	slice := it.Slice(s2)
	fmt.Println(slice)
	// Output: [aaaaaaa aaa aa a]
}
