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

func (g Chain[T]) Index() Chain2[int, T] {
	return Chain2[int, T]{seq: Index(g.seq)}
}

func (g Chain[T]) Reduce(reduceFunc ReduceFunc[T], initial T) T {
	return Reduce(g.seq, reduceFunc, initial)
}

func (g Chain[T]) Slice() []T {
	return Slice(g.seq)
}

// Mapable allows the operations to be chained via method calls and
// additionally T -> V and V -> T mapping can be added
type Mapable[T, V any] struct {
	seq  iter.Seq[T]
	none V
}

func NewMapable[T, V any](seq iter.Seq[T]) Mapable[T, V] {
	return Mapable[T, V]{
		seq: seq,
	}
}

func (g Mapable[T, V]) Filter(filterFunc FilterFunc[T]) Mapable[T, V] {
	return Mapable[T, V]{
		seq: Filter(g.seq, filterFunc),
	}
}

func (g Mapable[T, V]) Index() Chain2[int, T] {
	return Chain2[int, T]{seq: Index(g.seq)}
}

func (g Mapable[T, V]) Map(mapFunc MapFunc[T, V]) Mapable[V, T] {
	return Mapable[V, T]{
		seq: Map(g.seq, mapFunc),
	}
}

func (g Mapable[T, V]) Reduce(reduceFunc ReduceFunc[T], initial T) T {
	return Reduce(g.seq, reduceFunc, initial)
}

func (g Mapable[T, V]) Slice() []T {
	return Slice(g.seq)
}
