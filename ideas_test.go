package it_test

import (
	"fmt"
	"iter"
	"strconv"

	"github.com/gomoni/it"
)

type MapSeq2Func[T, K, V any] func(T) (K, V)

func MapSeq2[T, K, V any](seq iter.Seq[T], mapFunc MapSeq2Func[T, K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		for {
			t, ok := next()
			if !ok {
				return
			}
			k, v := mapFunc(t)
			if !yield(k, v) {
				return
			}
		}
	}
}

func withErrorFunc[T any, V error](t T) (T, error) {
	return t, nil
}

func WithError[T any](seq iter.Seq[T]) iter.Seq2[T, error] {
	return MapSeq2(seq, withErrorFunc[T, error])
}

func Example_idea_toseq2() {
	n := []string{"forty-two", "42"}
	s0 := it.From(n)
	s1 := WithError(s0) // this step is probably unnecessary and MapSeq2 would be used here
	s2 := it.Map2(s1, func(s string, _ error) (int, error) { return strconv.Atoi(s) })

	for value, error := range s2 {
		fmt.Println(value, error)
	}
	// Output:
	// 0 strconv.Atoi: parsing "forty-two": invalid syntax
	// 42 <nil>

}

func Example_idea_enumerable() {
	// instead of providing stuff like
	// type FilterIndexFunc[T any] func(T, int) bool
	// type Filter2IndexFunc
	// lets ignore the problem and simply solve that via a closure

	n := []string{"aa", "aaa", "aaaaaaa", "a"}

	// Map(enumerable) can be provided as a helper by it library
	type Indexed[T any] struct {
		index int
		value T
	}
	var idx int
	enumerable := func(s string) Indexed[string] {
		ret := Indexed[string]{index: idx, value: s}
		idx++
		return ret
	}
	res := it.NewMappableChain[string, Indexed[string]](it.From(n)).
		Map(enumerable).
		Filter(func(p Indexed[string]) bool { return p.index >= 2 }).
		Slice()
	fmt.Println(res)
	// Output: [{2 aaaaaaa} {3 a}]
}

type MapFuncError[T, V any] func(T) (V, error)

// Map calls a mapping function on each member of the sequence
// the consequence is an usage of iter.Seq2 - which avoids chaining of operations
func MapError[T, V any](s iter.Seq[T], mapFunc MapFuncError[T, V]) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		next, stop := iter.Pull(s)
		defer stop()

		for {
			t, ok := next()
			if !ok {
				return
			}
			v, err := mapFunc(t)
			if !yield(v, err) {
				return
			}
		}
	}
}

func Example_idea_errors() {
	// there are two areas
	// 1. simpler - helpers returning an error
	type FilterFuncError[T any] func(T) (bool, error)
	type ReduceFuncError[T any] func(T, T) (T, error)

	mapErrorFunc := func(s string) (int, error) {
		return strconv.Atoi(s)
	}

	// 2. even simpler due the Seq2 usage - the failible stuff can go as Seq2
	// so probably nothing for a simple chain

	n := []string{"forty-two", "42"}
	it0 := it.From(n)
	it1 := MapError(it0, mapErrorFunc)

	for v, err := range it1 {
		fmt.Println(v, err)
	}
	// Output:
	// 0 strconv.Atoi: parsing "forty-two": invalid syntax
	// 42 <nil>
}

type pusher struct {
	stack []string
}

func (y *pusher) push(s string) {
	y.stack = append(y.stack, s)
}

func (y pusher) seq() func(func(string) bool) {
	return it.From[string](y.stack)
}

func Example_break_da_chain() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}

	chain := it.NewChain(it.From(n)).
		Filter(func(s string) bool { return true })

	p := pusher{}
	for s := range chain.Seq() {
		p.push(s)
	}

	chain2 := it.NewChain[string](p.seq())
	slice := chain2.Slice()
	fmt.Println(slice)
	// Output: []
}
