package it_test

import (
	"fmt"
	"strings"

	"github.com/gomoni/it"
)

func ExampleFilter2() {
	m := map[string]int{"one": 0, "two": 1, "three": 2}

	s0 := it.From2(m)
	s1 := it.Filter2(s0, func(k string, v int) bool { return v >= 1 })
	for k, v := range it.Sorted2(s1, strings.Compare) {
		fmt.Println(k, v)
	}
	// Output:
	// three 2
	// two 1
}
