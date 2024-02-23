package it_test

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gomoni/it"
)

func ExampleFilter2() {
	m := map[string]int{"one": 0, "two": 1, "three": 2}

	s0 := it.From2(m)
	s1 := it.Filter2(s0, func(k string, v int) bool { return v >= 1 })
	s2 := it.Sort2(s1, func(slice []string) { slices.SortFunc(slice, strings.Compare) })
	for k, v := range s2 {
		fmt.Println(k, v)
	}
	// Output:
	// three 2
	// two 1
}

func ExampleMap2() {
	m := map[string]int{"one": 0, "two": 1, "three": 2}

	s0 := it.From2(m)
	s1 := it.Map2(s0, func(k string, v int) (int, string) { return v, k })
	s2 := it.Sort2(s1, slices.Sort)
	for k, v := range s2 {
		fmt.Println(k, v)
	}
	// Output:
	// 0 one
	// 1 two
	// 2 three
}

func ExampleKeys() {
	m := map[string]int{"one": 0, "two": 1, "three": 2}

	s0 := it.From2(m)
	s1 := it.Keys(s0)
	s2 := it.Sort(s1, slices.Sort)
	for s := range s2 {
		fmt.Println(s)
	}
	// Output:
	// one
	// three
	// two
}

func ExampleValues() {
	m := map[string]int{"one": 0, "two": 1, "three": 2}

	s0 := it.From2(m)
	s1 := it.Values(s0)
	s2 := it.Sort(s1, slices.Sort)
	for n := range s2 {
		fmt.Println(n)
	}
	// Output:
	// 0
	// 1
	// 2
}

func ExampleAsMap() {
	m := map[string]int{"one": 0, "two": 1, "three": 2}

	s0 := it.From2(m)
	s1 := it.Filter2(s0, func(_ string, v int) bool { return v == 2 })
	m2 := it.AsMap(s1)
	for k, v := range m2 {
		fmt.Println(k, v)
	}
	// Output:
	// three 2
}

func ExampleChain2() {
	m := map[string]int{"one": 0, "two": 1, "three": 2}

	m2 := it.NewChain2(it.From2(m)).
		Filter2(func(_ string, v int) bool { return v == 2 }).
		AsMap()
	for k, v := range m2 {
		fmt.Println(k, v)
	}
	// Output:
	// three 2
}
