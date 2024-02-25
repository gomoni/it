package it

import (
	"iter"
)

// Map2Func maps the K, V -> K2, V2
type Map2Func[K, V, K2, V2 any] func(K, V) (K2, V2)

// Filter2Func is a predicate for type K, V
type Filter2Func[K, V any] func(K, V) bool

func From2Slice[K comparable, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				break
			}
		}
	}
}

func Filter2[K, V any](seq iter.Seq2[K, V], filterFunc Filter2Func[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		next, stop := iter.Pull2(seq)
		defer stop()

		for {
			k, v, ok := next()
			if !ok {
				return
			}
			if shouldYield := filterFunc(k, v); !shouldYield {
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

func Map2[K, V, K2, V2 any](seq iter.Seq2[K, V], mapFunc Map2Func[K, V, K2, V2]) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		next, stop := iter.Pull2(seq)
		defer stop()

		for {
			k, v, ok := next()
			if !ok {
				return
			}
			k2, v2 := mapFunc(k, v)
			if !yield(k2, v2) {
				return
			}
		}
	}
}

func Keys[K any, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		next, stop := iter.Pull2(seq)
		defer stop()

		for {
			k, _, ok := next()
			if !ok {
				return
			}
			if !yield(k) {
				return
			}
		}
	}
}

func Values[K any, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		next, stop := iter.Pull2(seq)
		defer stop()

		for {
			_, v, ok := next()
			if !ok {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

func Sort2[K comparable, V any](seq iter.Seq2[K, V], sortFunc SortFunc[K]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {

		m := AsMap(seq)
		keys := make([]K, len(m))
		idx := 0
		for k := range m {
			keys[idx] = k
			idx++
		}
		sortFunc(keys)

		for idx := range len(m) {
			k := keys[idx]
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// AsMap converts the iter.Seq2 into map
func AsMap[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	ret := make(map[K]V, 1024)
	next, stop := iter.Pull2(seq)
	defer stop()

	for {
		k, v, ok := next()
		if !ok {
			return ret
		}
		ret[k] = v
	}
}
