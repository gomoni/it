package it

import (
	"iter"
	"slices"

	"github.com/gomoni/it/islices"
)

type Chain[T any] iter.Seq[T]

func NewChain[T any](seq iter.Seq[T]) Chain[T] {
	return Chain[T](seq)
}

func (ch Chain[T]) Seq() iter.Seq[T] {
	return iter.Seq[T](ch)
}

func (g Chain[T]) Filter(filterFunc islices.FilterFunc[T]) Chain[T] {
	return Chain[T](islices.Filter(g.Seq(), filterFunc))
}

func (g Chain[T]) Collect() []T {
	return slices.Collect(g.Seq())
}

type Mappable[T, V any] struct {
	seq  iter.Seq[T]
	none V
}

func NewMappable[T, V any](seq iter.Seq[T]) Mappable[T, V] {
	return Mappable[T, V]{
		seq: seq,
	}
}

func (g Mappable[T, V]) Seq() iter.Seq[T] {
	return g.seq
}

func (g Mappable[T, V]) Filter(filterFunc islices.FilterFunc[T]) Mappable[T, V] {
	return Mappable[T, V]{
		seq: islices.Filter(g.seq, filterFunc),
	}
}

func (g Mappable[T, V]) Map(mapFunc islices.MapFunc[T, V]) Mappable[V, T] {
	return Mappable[V, T]{
		seq: islices.Map(g.seq, mapFunc),
	}
}

func (g Mappable[T, V]) Collect() []T {
	return slices.Collect(g.Seq())
}
