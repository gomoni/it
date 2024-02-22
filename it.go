package it

import (
	"iter"
	"slices"
)

// FilterFunc is a predicate for type T
type FilterFunc[T any] func(T) bool

// MapFunc maps the T -> V
type MapFunc[T, V any] func(T) V

// Reduce func combines two arguments into one
type ReduceFunc[T any] func(T, T) T

// From converts the slice into iter.Seq
// TODO: better names like FromSlice, FromMap, FromChannel?
func From[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, i := range slice {
			if !yield(i) {
				break
			}
		}
	}
}

// SimpleFilter is a simple wrapper of one Filter operation on a slice []T.
func SimpleFilter[T any](slice []T, filterFunc FilterFunc[T]) []T {
    oseq := From(slice)
    fseq := Filter(oseq, filterFunc)
    return Slice(fseq)
}

// SimpleFilters is a simple wrapper over several consecutive Filter operations ona slice []T.
func SimpleFilters[T any](slice []T, filterFuncs ...FilterFunc[T]) []T {
    var rv []T
    rv = slice
    for _, filterFunc := range filterFuncs {
        rv = SimpleFilter(rv, filterFunc)
    }
    return rv
}

// Filter yields only values for which filterFunc returns true
func Filter[T any](s iter.Seq[T], filterFunc FilterFunc[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(s)
		defer stop()

		for {
			t, ok := next()
			if !ok {
				return
			}
			if shouldYield := filterFunc(t); !shouldYield {
				continue
			}
			if !yield(t) {
				return
			}
		}
	}
}

// Map calls a mapping function on each member of the sequence
func Map[T, V any](s iter.Seq[T], mapFunc MapFunc[T, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		next, stop := iter.Pull(s)
		defer stop()

		for {
			t, ok := next()
			if !ok {
				return
			}
			v := mapFunc(t)
			if !yield(v) {
				return
			}
		}
	}
}

func Reduce[T any](s iter.Seq[T], reduceFunc ReduceFunc[T], initial T) T {
	ret := initial
	next, stop := iter.Pull(s)
	defer stop()

	for {
		t, ok := next()
		if !ok {
			return ret
		}
		ret = reduceFunc(ret, t)
	}
}

func Reverse[T any](s iter.Seq[T]) iter.Seq[T] {
	slice := Slice(s)

	return func(yield func(T) bool) {
		for idx := len(slice) - 1; idx != -1; idx-- {
			if !yield(slice[idx]) {
				break
			}
		}
	}
}

func Sorted[T any](s iter.Seq[T], sortFunc SortFunc[T]) iter.Seq[T] {
	slice := Slice(s)
	slices.SortFunc(slice, sortFunc)
	return From(slice)
}

// Slice converts the sequence into slice
// TODO: rename to AsSlice like AsMap?
// TODO: provide IntoScile(slice []T, seq iter.Seq[T])?
func Slice[T any](seq iter.Seq[T]) []T {
	ret := make([]T, 0, 1024)
	next, stop := iter.Pull(seq)
	defer stop()

	for {
		t, ok := next()
		if !ok {
			return ret
		}
		ret = append(ret, t)
	}
}
