package it_test

import (
	"testing"

	"github.com/gomoni/it"
)

var in1M []int

func init() {
	in1M = make([]int, 1024*1024)
	for idx := range in1M {
		in1M[idx] = idx
	}
}

func rangeFilter(in []int, filterFunc it.FilterFunc[int]) []int {
	ret := make([]int, 0, len(in))
	for _, i := range in {
		if filterFunc(i) {
			ret = append(ret, i)
		}
	}
	return ret
}

func BenchmarkRangeFilter(b *testing.B) {
	for range b.N {
		rangeFilter(in1M, func(i int) bool { return i%2 == 0 })
	}
}

func BenchmarkItFilterSlice(b *testing.B) {
	for range b.N {
		seq0 := it.Filter(it.FromSlice(in1M), func(i int) bool { return i%2 == 0 })
		it.AsSlice(seq0)
	}
}

func BenchmarkItFilterFor(b *testing.B) {
	for range b.N {
		seq0 := it.Filter(it.FromSlice(in1M), func(i int) bool { return i%2 == 0 })
		cnt := 0
		for value := range seq0 {
			cnt += value
		}
	}
}
