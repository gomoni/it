package it_test

import (
	"slices"
	"testing"

	"github.com/gomoni/it/islices"
)

const size = 1024 * 1024

var in1M []int

func init() {
	in1M = make([]int, size)
	for idx := range in1M {
		in1M[idx] = idx
	}
}

// BenchmarkRange benchmarks the for range loop over a slice
func BenchmarkRange(b *testing.B) {
	for range b.N {
		cnt := 0
		for _, value := range in1M {
			cnt += value
		}
	}
}

// BenchmarkRangeAll benchmarks the slices.All
func BenchmarkRangeAll(b *testing.B) {
	for range b.N {
		cnt := 0
		for _, value := range slices.All(in1M) {
			cnt += value
		}
	}
}

// BenchmarkRangeValues benchmarks the slices.Values
func BenchmarkRangeValues(b *testing.B) {
	for range b.N {
		cnt := 0
		for value := range slices.Values(in1M) {
			cnt += value
		}
	}
}

// BenchmarkRangeAll benchmarks a range loop skipping the odd numbers
func BenchmarkRangeEven(b *testing.B) {
	for range b.N {
		cnt := 0
		for _, value := range in1M {
			if value%2 != 0 {
				continue
			}
			cnt += value
		}
	}
}

// BenchmarkRangeFilterEven uses a Filter method on a sequence to do the filtering
func BenchmarkRangeValuesFilterEven(b *testing.B) {
	for range b.N {
		cnt := 0
		all := slices.Values(in1M)
		evens := islices.Filter(all, func(value int) bool {
			return value%2 == 0
		})
		for value := range evens {
			cnt += value
		}
	}
}
