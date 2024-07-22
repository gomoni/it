package imaps_test

import (
	"fmt"
	"maps"

	imaps "github.com/gomoni/it/imaps"
)

func ExampleFilter() {
	m := map[string]int{
		"bambino": 1,
		"junior":  11,
		"senior":  22,
	}
	s0 := maps.All(m)
	s1 := imaps.Filter(s0, func(_ string, v int) bool { return v >= 18 })
	result := maps.Collect(s1)
	fmt.Println(result)
	// Output: map[senior:22]
}

func ExampleMap() {
	m := map[string]int{
		"bambino": 1,
		"junior":  11,
		"senior":  22,
	}
	s0 := maps.All(m)
	s1 := imaps.Map(s0, func(s string, _ int) (string, int) { return s, len(s) })
	result := maps.Collect(s1)
	fmt.Println(result)
	// Output: map[bambino:7 junior:6 senior:6]
}
