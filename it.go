package it

import (
	"iter"
)

// FilterFunc is a predicate for type T
type FilterFunc[T any] func(T) bool

// MapFunc maps the T -> V
type MapFunc[T, V any] func(T) V

// Reduce func combines two arguments into one
type ReduceFunc[T any] func(T, T) T

// Sort func sorts the given slice
type SortFunc[T any] func([]T)

// MapSeq2Func allows a mapping of single T into two other types
// A typical example is to return a T with an error in case mapping fails
type MapSeq2Func[T, K, V any] func(T) (K, V)

type MapErrorFunc[T, V any] func(T) (V, error)

// FromSlice converts the slice into iter.Seq
func FromSlice[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, i := range slice {
			if !yield(i) {
				break
			}
		}
	}
}

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

// IndexFrom returns sequence returning two values index, T starting from initial.
// This is compatible with a standard  range over a slice
//
//	for index, value := range IndexFrom(sequence, 42) {}
func IndexFrom[T any](seq iter.Seq[T], initial int) iter.Seq2[int, T] {
	index := initial
	return func(yield func(int, T) bool) {
		for v := range seq {
			if !yield(index, v) {
				return
			}
			index++
		}
	}
}

// Index returns sequence returning two values index, T starting from zero.
// See IndexFrom for details
func Index[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	return IndexFrom(seq, 0)
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

// MapSeq2 maps the single value of type T into a sequence of two values of types K, V
// and a typical use is return a mapped type together with an error
func MapSeq2[T, K, V any](seq iter.Seq[T], mapFunc MapSeq2Func[T, K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for v := range seq {
			k, v := mapFunc(v)
			if !yield(k, v) {
				return
			}
		}
	}
}

func mapErrorFunc[T, K any](mapFunc MapErrorFunc[T, K]) MapSeq2Func[T, K, error] {
	return func(value T) (K, error) {
		return mapFunc(value)
	}
}

func MapError[T, V any](seq iter.Seq[T], mapFunc MapErrorFunc[T, V]) iter.Seq2[V, error] {
	return MapSeq2[T, V, error](seq, mapErrorFunc(mapFunc))
}

func Reduce[T any](s iter.Seq[T], reduceFunc ReduceFunc[T], initial T) T {
	ret := initial

	for v := range s {
		ret = reduceFunc(ret, v)
	}

	return ret
}

func Reverse[T any](s iter.Seq[T]) iter.Seq[T] {
	slice := AsSlice(s)

	return func(yield func(T) bool) {
		for idx := len(slice) - 1; idx != -1; idx-- {
			if !yield(slice[idx]) {
				break
			}
		}
	}
}

func Sort[T any](s iter.Seq[T], sortFunc SortFunc[T]) iter.Seq[T] {
	slice := AsSlice(s)
	sortFunc(slice)
	return FromSlice(slice)
}

// AsSlice converts the sequence into slice
// TODO: rename to AsSlice like AsMap?
// TODO: provide IntoScile(slice []T, seq iter.Seq[T])?
func AsSlice[T any](seq iter.Seq[T]) []T {
	ret := make([]T, 0, 1024)

	for v := range seq {
		ret = append(ret, v)
	}

	return ret
}
