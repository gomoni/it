// Package maps defines various iterators useful with value tuples of any type.

package imaps

import (
	"iter"
)

// Map2Func maps the K, V -> K2, V2
type Map2Func[K, V, K2, V2 any] func(K, V) (K2, V2)

// Filter2Func is a predicate for type K, V
type Filter2Func[K, V any] func(K, V) bool

func Filter[K, V any](s2 iter.Seq2[K, V], filterFunc Filter2Func[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s2 {
			if shouldYield := filterFunc(k, v); !shouldYield {
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

func Map[K, V, K2, V2 any](s2 iter.Seq2[K, V], mapFunc Map2Func[K, V, K2, V2]) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s2 {
			k2, v2 := mapFunc(k, v)
			if !yield(k2, v2) {
				return
			}
		}
	}
}
