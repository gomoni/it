package it_test

import (
	"fmt"
	"slices"
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

func ExampleIndexFrom() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	s0 := it.From(n)
	for index, value := range it.IndexFrom(s0, 42) {
		fmt.Println(index, value)
	}
	// Output:
	// 42 aa
	// 43 aaa
	// 44 aaaaaaa
	// 45 a
}

func ExampleIndex() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	s0 := it.From(n)
	for index, value := range it.Index(s0) {
		fmt.Println(index, value)
	}
	// Output:
	// 0 aa
	// 1 aaa
	// 2 aaaaaaa
	// 3 a
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

func ExampleMapSeq2() {
	n := []string{"forty-two", "42"}
	s0 := it.From(n)
	s1 := it.MapSeq2(s0, func(s string) (int, error) { return strconv.Atoi(s) })
	for value, error := range s1 {
		fmt.Println(value, error)
	}
	// Output:
	// 0 strconv.Atoi: parsing "forty-two": invalid syntax
	// 42 <nil>
}

func ExampleMapError() {
	n := []string{"forty-two", "42"}
	s0 := it.From(n)
	s1 := it.MapError(s0, strconv.Atoi)
	for value, error := range s1 {
		fmt.Println(value, error)
	}
	// Output:
	// 0 strconv.Atoi: parsing "forty-two": invalid syntax
	// 42 <nil>
}

func ExampleSort() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	s0 := it.From(n)
	s1 := it.Sort(s0, func(slice []string) { slices.SortFunc(slice, strings.Compare) })
	slice := it.Slice(s1)
	fmt.Println(slice)
	// Output: [a aa aaa aaaaaaa]
}

func ExampleReverse() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	s0 := it.From(n)
	s1 := it.Sort(s0, func(slice []string) { slices.SortFunc(slice, strings.Compare) })
	s2 := it.Reverse(s1)
	slice := it.Slice(s2)
	fmt.Println(slice)
	// Output: [aaaaaaa aaa aa a]
}
