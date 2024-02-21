package it

import "iter"

// Chain allows the sequence operations to be chained via method calls
// works with a single type
type Chain[T any] struct {
	seq iter.Seq[T]
}

func NewChain[T any](seq iter.Seq[T]) Chain[T] {
	return Chain[T]{seq: seq}
}

func (g Chain[T]) Seq() iter.Seq[T] {
	return g.seq
}

func (g Chain[T]) Filter(filterFunc FilterFunc[T]) Chain[T] {
	return Chain[T]{seq: Filter(g.seq, filterFunc)}
}

func (g Chain[T]) Reduce(reduceFunc ReduceFunc[T], initial T) T {
	return Reduce(g.seq, reduceFunc, initial)
}

func (g Chain[T]) Slice() []T {
	return Slice(g.seq)
}

// Mappable allows the operations to be chained via method calls and
// additionally T -> V and V -> T mapping can be added
// XXX: naming - the 2 suffix leads one to think this does work with maps :shrug:
type Mappable[T, V any] struct {
	seq  iter.Seq[T]
	none V
}

func NewMappable[T, V any](seq iter.Seq[T]) Mappable[T, V] {
	return Mappable[T, V]{
		seq: seq,
	}
}

func (g Mappable[T, V]) Filter(filterFunc FilterFunc[T]) Mappable[T, V] {
	return Mappable[T, V]{
		seq: Filter(g.seq, filterFunc),
	}
}

func (g Mappable[T, V]) Map(mapFunc MapFunc[T, V]) Mappable[V, T] {
	return Mappable[V, T]{
		seq: Map(g.seq, mapFunc),
	}
}

func (g Mappable[T, V]) Reduce(reduceFunc ReduceFunc[T], initial T) T {
	return Reduce(g.seq, reduceFunc, initial)
}

func (g Mappable[T, V]) Slice() []T {
	return Slice(g.seq)
}
