// Package islices defines various iterators useful with slices of any type.

package islices

import "iter"

// FilterFunc is a predicate for type T
type FilterFunc[T any] func(T) bool

// MapFunc maps the T -> V
type MapFunc[T, V any] func(T) V

// Filter yields only values for which filterFunc returns true
func Filter[T any](s iter.Seq[T], filterFunc FilterFunc[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if shouldYield := filterFunc(v); !shouldYield {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Map calls a mapping function on each member of the sequence
func Map[T, V any](s iter.Seq[T], mapFunc MapFunc[T, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			v := mapFunc(v)
			if !yield(v) {
				return
			}
		}
	}
}
