package islices_test

import (
	"fmt"
	"slices"
	"strconv"

	islices "github.com/gomoni/it/islices"
)

func ExampleFilter() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	s0 := slices.Values(n)
	s1 := islices.Filter(s0, func(s string) bool { return len(s) >= 2 })
	slice := slices.Collect(s1)
	fmt.Println(slice)
	// Output: [aa aaa aaaaaaa]
}

func ExampleMap() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	// maps string->int->float32
	s0 := slices.Values(n)
	s1 := islices.Map(s0, func(s string) int { return len(s) })
	s2 := islices.Map(s1, func(i int) float32 { return float32(i) })
	s3 := islices.Map(s2, func(f float32) string { return strconv.FormatFloat(float64(f), 'E', 4, 32) })
	slice := slices.Collect(s3)
	fmt.Println(slice)
	// Output: [2.0000E+00 3.0000E+00 7.0000E+00 1.0000E+00]
}
