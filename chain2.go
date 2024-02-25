package it

import "iter"

// Chain2 allows the sequence operations to be chained via method calls
// works with a single type
type Chain2[K, V any] struct {
	seq iter.Seq2[K, V]
}

func NewChain2[K, V any](seq iter.Seq2[K, V]) Chain2[K, V] {
	return Chain2[K, V]{seq: seq}
}

func (g Chain2[K, V]) Seq2() iter.Seq2[K, V] {
	return g.seq
}

func (g Chain2[K, V]) Filter2(filterFunc Filter2Func[K, V]) Chain2[K, V] {
	return Chain2[K, V]{seq: Filter2(g.seq, filterFunc)}
}

func (g Chain2[K, V]) Keys() Chain[K] {
	return Chain[K]{seq: Keys(g.seq)}
}

func (g Chain2[K, V]) Values() Chain[V] {
	return Chain[V]{seq: Values(g.seq)}
}
