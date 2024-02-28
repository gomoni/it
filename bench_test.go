package it_test

import (
	"iter"
	"testing"

	"github.com/gomoni/it"
)

const size = 1024 * 1024

var in1M []int

func init() {
	in1M = make([]int, size)
	for idx := range in1M {
		in1M[idx] = idx
	}
}

func push(size int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for value := range size {
			if !yield(value) {
				break
			}
		}
	}
}

// BenchmarkRange benchmarks the iterator pushing from a slice
// baseline: ~2350 ops
func BenchmarkRange(b *testing.B) {
	for range b.N {
		cnt := 0
		for value := range in1M {
			cnt += value
		}
	}
}

// BenchmarkPushIterator benchmarks the iterator pushing from a slice
// benchmarks the overhead of a pull iterator over plain range version
// overhead over BenchmarkRange is less than 10%
func BenchmarkPushIterator(b *testing.B) {
	for range b.N {
		cnt := 0
		for value := range push(size) {
			cnt += value
		}
	}
}

func pull(seq iter.Seq[int]) iter.Seq[int] {
	return func(yield func(int) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		for {
			t, ok := next()
			if !ok || !yield(t) {
				return
			}
		}
	}
}

// BenchmarkPullIterator tests slice -> push -> pull
func BenchmarkPullIterator(b *testing.B) {
	for range b.N {
		cnt := 0
		for value := range pull(push(size)) {
			cnt += value
		}
	}
}

func rangeFilterSlice(in []int, filterFunc it.FilterFunc[int]) []int {
	ret := make([]int, 0, len(in))
	for _, i := range in {
		if filterFunc(i) {
			ret = append(ret, i)
		}
	}
	return ret
}

// BenchmarkRangeFilter benchmarks the for index, value := range loop with a filter - this is a baseline
func BenchmarkRangeFilterSlice(b *testing.B) {
	for range b.N {
		rangeFilterSlice(in1M, func(i int) bool { return i%2 == 0 })
	}
}

// Benchmark filters push -> pull -> pull operation
func BenchmarkItFilterSlice(b *testing.B) {
	for range b.N {
		seq0 := it.Filter(it.FromSlice(in1M), func(i int) bool { return i%2 == 0 })
		it.AsSlice(seq0)
	}
}

// Benchmark filters push -> pull -> pull but does not allocate
func BenchmarkItFilterFor(b *testing.B) {
	for range b.N {
		seq0 := it.Filter(it.FromSlice(in1M), func(i int) bool { return i%2 == 0 })
		cnt := 0
		for value := range seq0 {
			cnt += value
		}
	}
}
