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
