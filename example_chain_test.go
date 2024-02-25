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

func ExampleChain_Index() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	res := it.NewChain(it.From(n)).
		Index().
		Filter2(func(index int, _ string) bool { return index <= 1 }).
		Values().
		Slice()
	fmt.Println(res)
	// Output: [aa aaa]
}

func ExampleChain_Reduce() {
	m := []int{1, 2, 3, 4, 5, 6, 7}
	count := it.NewChain(it.From(m)).
		Reduce(func(a, _ int) int { return a + 1 }, 0)
	fmt.Println(count)
	// Output: 7
}

func Example_readme_chain() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	slice := it.NewMapable[string, int](it.From(n)).
		Map(func(s string) int { return len(s) }).
		Index().
		Filter2(func(index int, _ int) bool { return index <= 1 }).
		Values().
		Slice()
	fmt.Println(slice)
	// Output: [2 3]
}

func Example_readme_plain() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	seq0 := it.From(n)
	seq1 := it.Map(seq0, func(s string) int { return len(s) })
	seq2 := it.Index(seq1)
	seq3 := it.Filter2(seq2, func(index int, _ int) bool { return index <= 1 })
	seq4 := it.Values(seq3)
	slice := it.Slice(seq4)
	fmt.Println(slice)
	// Output: [2 3]
}

func ExampleMapable_Map() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	// maps string->int
	res := it.NewMapable[string, int](it.From(n)).
		Map(func(s string) int { return len(s) }).
		Filter(func(i int) bool { return i >= 2 }).
		Slice()
	fmt.Println(res)
	// Output: [2 3 7]
}

func ExampleMapable_MapError() {
	n := []string{"forty-two", "42"}
	c := it.NewMapable[string, int](it.From(n)).
		MapError(strconv.Atoi)
	for value, error := range c.Seq2() {
		fmt.Println(value, error)
	}
	// Output:
	// 0 strconv.Atoi: parsing "forty-two": invalid syntax
	// 42 <nil>
}

func ExampleMapable_Map_second() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	// maps string->int and int->string
	res := it.NewMapable[string, int](it.From(n)).
		Map(func(s string) int { return len(s) }).
		Filter(func(i int) bool { return i >= 2 }).
		Map(func(i int) string { return "string(" + strconv.Itoa(i) + ")" }).
		Slice()
	fmt.Println(res)
	// Output: [string(2) string(3) string(7)]
}
