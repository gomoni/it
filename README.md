# it: experimental iterator utilities for Go

# The problem

Go lacks an ergonomically composable idiomatic iterator utility library. Package `it` builds
on top of a (rangefunc experiment)[https://go.dev/wiki/RangefuncExperiment] for
go 1.22.

The design goals are

 * Minimal library providing just enough.
 * Idiomatic Go: that means use generic functions and use `rangefunc`
   experiment under the hood.
 * Type safe: so use generics everywhere.
 * Provide a composable API where practical
 * Supports Go builtin types: slices/maps/channels

Non goals are

 * Support of every iterator utility only because it exists in lo/lodash/Rust/Haskell/Scala/Python
 * Every possible permutation of primitives provides by the `it` itself.

# Usage

> Don't forget to install go 1.22 and `export GOEXPERIMENT=rangefunc`

The `it` provides methods that can be chained together. These are a better show
case to developers.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
slice := it.NewMapable[string, int](it.From(n)).
	Map(func(s string) int { return len(s) }).
	Index().
	Filter2(func(index int, _ int) bool { return index <= 1 }).
	Values().
	Slice()
fmt.Println(slice)
// Output: [2 3]
```

All this can be done using simple functions and an explicit sequence passing.
Go compiler will catch the unused variables and type mismatches helping the
developer in this form too.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
seq0 := it.From(n)
seq1 := it.Map(seq0, func(s string) int { return len(s) })
seq2 := it.Index(seq1)
seq3 := it.Filter2(seq2, func(index int, _ int) bool { return index <= 1 })
seq4 := it.Values(seq3)
slice := it.Slice(seq4)
fmt.Println(slice)
// Output: [2 3]
```

# Examples

## Filtering

In order to limit the sequence, use a `Filter`.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
res := it.NewChain(it.From(n)).
	Filter(func(s string) bool { return len(s) >= 2 }).
	Filter(func(s string) bool { return len(s) >= 3 }).
	Slice()
fmt.Println(res)
// Output: [aaa aaaaaaa]
```

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
s0 := it.From(n)
s1 := it.Filter(s0, func(s string) bool { return len(s) >= 2 })
slice := it.Slice(s1)
fmt.Println(slice)
// Output: [aa aaa aaaaaaa]
```

## Indexing

Every Go developer is familiar with iterating through a slice and two variable range form.

```go
for index, value := range slice {}
```

The `it` does implement the `Index`/`IndexFrom` functions, which adds the index into the sequence.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
s0 := it.From(n)
for index, value := range it.Index(s0) {
	fmt.Println(index, value)
}
// Output:
// 0 aa
// 1 aaa
// 2 aaaaaaa
// 3 a
```

The way how this is implemented is that `it` returns `iter.Seq2[int, T]`. This
type uses functions suffixed by `2`. That means `Filter2` is used in this
example. The `Values` functions returns the second value later on.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
res := it.NewChain(it.From(n)).
	Index().
	Filter2(func(index int, _ string) bool { return index <= 1 }).
	Values().
	Slice()
fmt.Println(res)
// Output: [aa aaa]
```

## Map

Map transforms one type into another. This is easy to do in Go as a simple
function. Much harder to do via method. The `Mapable` allows the developer to
use a `Map` in a method chain.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
// maps string->int and int->string
res := it.NewMapable[string, int](it.From(n)).
	Map(func(s string) int { return len(s) }).
	Filter(func(i int) bool { return i >= 2 }).
	Map(func(i int) string { return "string(" + strconv.Itoa(i) + ")" }).
	Slice()
fmt.Println(res)
// Output: [string(2) string(3) string(7)]
```

There are only two drawbacks

 1. Developer has to specify type parameters in advance
 2. It supports only two types - supporting more would lead to confusing API

However good old functions have no such limitation and can be used instead.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
// maps string->int->float32
s0 := it.From(n)
s1 := it.Map(s0, func(s string) int { return len(s) })
s2 := it.Map(s1, func(i int) float32 { return float32(i) })
s3 := it.Map(s2, func(f float32) string { return strconv.FormatFloat(float64(f), 'E', 4, 32) })
slice := it.Slice(s3)
fmt.Println(slice)
// Output: [2.0000E+00 3.0000E+00 7.0000E+00 1.0000E+00]
```


## Map with errors

Sometimes there is no 1:1 transformation between `T` and `V` and the mapping
can fail. The `it` provides a mapping function from `iter.Seq[T]` into
`iter.Seq2[K, V]`, which can solve this problem.

```go
n := []string{"forty-two", "42"}
s0 := it.From(n)
s1 := it.MapSeq2(s0, strconv.Atoi)
for value, error := range s1 {
	fmt.Println(value, error)
}
// Output:
// 0 strconv.Atoi: parsing "forty-two": invalid syntax
// 42 <nil>
```

Since this is a very common operation in Go the specialised `MapError` function
exists as an exception to the permutation rule above.

```go
n := []string{"forty-two", "42"}
s0 := it.From(n)
s1 := it.MapError(s0, strconv.Atoi)
for value, error := range s1 {
	fmt.Println(value, error)
}
// Output:
// 0 strconv.Atoi: parsing "forty-two": invalid syntax
// 42 <nil>
```

And can be done inside a method chain too.

```go
n := []string{"forty-two", "42"}
c := it.NewMapable[string, int](it.From(n)).
	MapError(strconv.Atoi)
for value, error := range c.Seq2() {
	fmt.Println(value, error)
}
// Output:
// 0 strconv.Atoi: parsing "forty-two": invalid syntax
// 42 <nil>
```

## Reduce

Reduce is a common functional operation, except it returns a single value. It
allows the developer to implement operation `Count`.

```go
m := []int{1, 2, 3, 4, 5, 6, 7}
count := it.NewChain(it.From(m)).
	Reduce(func(a, _ int) int { return a + 1 }, 0)
fmt.Println(count)
// Output: 7
```

## Sort

All other operations can work on a single item at the time. Not sort - it first
pulls all items to the slice, sorts them and then pushes the values to the iterator.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
s0 := it.From(n)
s1 := it.Sort(s0, func(slice []string) { slices.SortFunc(slice, strings.Compare) })
slice := it.Slice(s1)
fmt.Println(slice)
// Output: [a aa aaa aaaaaaa]
```

It accepts `type SortFunc[T any] func([]T)` instead of a `less` function. This
allows the developer to specify exactly _how_ the sequence should be
sorted. For example use a `slices.SortStableFunc` to get a stable sort.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
s0 := it.From(n)
s1 := it.Sort(s0, func(slice []string) { slices.SortStableFunc(slice, strings.Compare) })
slice := it.Slice(s1)
fmt.Println(slice)
// Output: [a aa aaa aaaaaaa]
```

## Reverse

Simply iterate backward - it must collect the slice first.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
s0 := it.From(n)
s1 := it.Sort(s0, func(slice []string) { slices.SortFunc(slice, strings.Compare) })
s2 := it.Reverse(s1)
slice := it.Slice(s2)
fmt.Println(slice)
// Output: [aaaaaaa aaa aa a]
```

## iter.Seq2[K, V] and Chain2

Most operations does have the alternative working on `iter.Seq2[K, V]`. In
order to distinguish between the types, all functions and structs has a suffix `2`, so
it is clear if method works with a single value or a pair.

Filtering

```go
m := map[string]int{"one": 0, "two": 1, "three": 2}

s0 := it.From2(m)
s1 := it.Filter2(s0, func(k string, v int) bool { return v >= 1 })
s2 := it.Sort2(s1, func(slice []string) { slices.SortFunc(slice, strings.Compare) })
for k, v := range s2 {
	fmt.Println(k, v)
}
// Output:
// three 2
// two 1
```

Note the sort is mandatory - the order of a range loop was changed from time to time.

```go
m := map[string]int{"one": 0, "two": 1, "three": 2}

s0 := it.From2(m)
s1 := it.Map2(s0, func(k string, v int) (int, string) { return v, k })
s2 := it.Sort2(s1, slices.Sort)
for k, v := range s2 {
	fmt.Println(k, v)
}
// Output:
// 0 one
// 1 two
// 2 three
```

Sometimes developer only one of `K`, `V`, so `Keys` and `Values` do not
have `2` suffix as they return `itre.Seq`.

```go
m := map[string]int{"one": 0, "two": 1, "three": 2}

s0 := it.From2(m)
s1 := it.Keys(s0)
s2 := it.Sort(s1, slices.Sort)
for s := range s2 {
	fmt.Println(s)
}
// Output:
// one
// three
// two
```

```go
m := map[string]int{"one": 0, "two": 1, "three": 2}

s0 := it.From2(m)
s1 := it.Values(s0)
s2 := it.Sort(s1, slices.Sort)
for n := range s2 {
	fmt.Println(n)
}
// Output:
// 0
// 1
// 2
```

And the developer can get values back as a map - note the `K` must be
`comparable` otherwise the type system will not allow one to construct a map. This
is the reason `Chain2` does not have `AsMap` method. Doing so would impose the
constraint to both `K` and `V` and limiting the usability of a `Chain2`.

```
invalid map key type K (missing comparable constraint)
```

```go
m := map[string]int{"one": 0, "two": 1, "three": 2}

s0 := it.From2(m)
s1 := it.Filter2(s0, func(_ string, v int) bool { return v == 2 })
m2 := it.AsMap(s1)
for k, v := range m2 {
	fmt.Println(k, v)
}
// Output:
// three 2
```

All operations above can be chained.

```go
m := map[string]int{"one": 0, "two": 1, "three": 2}

chain2 := it.NewChain2(it.From2(m)).
	Filter2(func(_ string, v int) bool { return v == 2 })
for k, v := range chain2.Seq2() {
	fmt.Println(k, v)
}
// Output:
// three 2
```

# Ideas

Some crazy and not so crazy ideas to expolse

## break the chain

One of the coolest (keep in mind it was a late night one) ideas may be breaking
the chain. The prototype exists in ideas_test.go, just not sure if it _is_
actually a good idea. It is definitely doable and possible in Go.

```go
package it_test

import (
	"fmt"

	"github.com/gomoni/it"
)

type pusher struct {
	stack chan string
}

func (y *pusher) push(s string) {
	y.stack <- s
}

func (y pusher) seq() func(func(string) bool) {
	return func(yield func(string) bool) {
		for {
			select {
			case s, open := <-y.stack:
				if !open || !yield(s) {
					return
				}
			}
		}
	}
}

func (y pusher) wait() {
	<-y.stack
}

func Example_break_da_chain() {
	n := []string{"aa", "aaa", "aaaaaaa", "a"}

	// create a method chain
	chain := it.NewChain(it.From(n)).
		Filter(func(s string) bool { return true })

	// break it - with some syntax sugar
	p := pusher{stack: make(chan string)}
	defer p.wait()
	go func() {
		defer close(p.stack)
		for s := range chain.Seq() {
			p.push(s)
		}
	}()

	// continue here
	chain2 := it.NewChain(p.seq()).
		Filter(func(s string) bool { return len(s) > 2 })
	slice := chain2.Slice()
	fmt.Println(slice)
	// Output: [aaa aaaaaaa]
}
```


## make iterations context aware?????

 Some options

 1. don't do that
 1. provide `FilterContext` et all
 1. `it/itctx` package

# Other libraries

Inspiration and other cool projects.

## lo

https://pkg.go.dev/github.com/samber/lo

Is the most used lodash-style library for Go.

Pros

 * most favorite
 * type safe due usage of generics
 * one can iterate over slices, maps or channels
 * helper function on everything

Cons

 * operate on top of simple functions
 * each pass allocates new slice or map
 * methods can't be chanined together
 * every callback got mandatory int argument

An example is `FilterMap` function. It can be implemented as `Filter` and
`Map`, yet it's not easy to do in `lo`

https://pkg.go.dev/github.com/samber/lo#pkg-functions

## gubrak

https://github.com/novalagung/gubrak

Less popular, provides nicer looking API than lo. Implements own iterator type,
so methods can be arbitrary chained. Harder to use due prevalent `interface{}` usage.

Pros

 * nicer API
 * methods can be chained
 * provides own iterator type

Cons

 * not type safe due usage of `interface{}`
 * API actually harder to use

## Other languages

 * https://docs.python.org/3/library/itertools.html
 * https://doc.rust-lang.org/book/ch13-02-iterators.html
 * https://hackage.haskell.org/package/base-4.19.0.0/docs/Data-List.html#g:2

Documentation generated using https://github.com/dave/rebecca

```
go install github.com/dave/rebecca/cmd/becca@latest
becca -package github.com/gomoni/it
```
