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

{{ "Example_readme_chain" | example }}

All this can be done using simple functions and an explicit sequence passing.
Go compiler will catch the unused variables and type mismatches helping the
developer in this form too.

{{ "Example_readme_plain" | example }}

# Examples

## Filtering

In order to limit the sequence, use a `Filter`.

{{ "ExampleChain_Filter" | example }}

{{ "ExampleFilter" | example }}

## Indexing

Every Go developer is familiar with iterating through a slice and two variable range form.

```go
for index, value := range slice {}
```

The `it` does implement the `Index`/`IndexFrom` functions, which adds the index into the sequence.

{{ "ExampleIndex" | example }}

The way how this is implemented is that `it` returns `iter.Seq2[int, T]`. This
type uses functions suffixed by `2`. That means `Filter2` is used in this
example. The `Values` functions returns the second value later on.

{{ "ExampleChain_Index" | example }}

## Map

Map transforms one type into another. This is easy to do in Go as a simple
function. Much harder to do via method. The `Mapable` allows the developer to
use a `Map` in a method chain.

{{ "ExampleMapable_Map_second" | example }}

There are only two drawbacks

 1. Developer has to specify type parameters in advance
 2. It supports only two types - supporting more would lead to confusing API

However good old functions have no such limitation and can be used instead.

{{ "ExampleMap" | example }}


## Map with errors

Sometimes there is no 1:1 transformation between `T` and `V` and the mapping
can fail. The `it` provides a mapping function from `iter.Seq[T]` into
`iter.Seq2[K, V]`, which can solve this problem.

{{ "ExampleMapSeq2" | example }}

Since this is a very common operation in Go the specialised `MapError` function
exists as an exception to the permutation rule above.

{{ "ExampleMapError" | example }}

And can be done inside a method chain too.

{{ "ExampleMapable_MapError" | example }}

## Reduce

Reduce is a common functional operation, except it returns a single value. It
allows the developer to implement operation `Count`.

{{ "ExampleChain_Reduce" | example }}

## Sort

All other operations can work on a single item at the time. Not sort - it first
pulls all items to the slice, sorts them and then pushes the values to the iterator.

{{ "ExampleSort" | example }}

It accepts `type SortFunc[T any] func([]T)` instead of a `less` function. This
allows the developer to specify exactly _how_ the sequence should be
sorted. For example use a `slices.SortStableFunc` to get a stable sort.

{{ "ExampleSort_stable" | example }}

## Reverse

Simply iterate backward - it must collect the slice first.

{{ "ExampleReverse" | example }}

## iter.Seq2[K, V] and Chain2

Most operations does have the alternative working on `iter.Seq2[K, V]`. In
order to distinguish between the types, all functions and structs has a suffix `2`, so
it is clear if method works with a single value or a pair.

Filtering

{{ "ExampleFilter2" | example }}

Note the sort is mandatory - the order of a range loop was changed from time to time.

{{ "ExampleMap2" | example }}

Sometimes developer only one of `K`, `V`, so `Keys` and `Values` do not
have `2` suffix as they return `itre.Seq`.

{{ "ExampleKeys" | example }}

{{ "ExampleValues" | example }}

And the developer can get values back as a map - note the `K` must be
`comparable` otherwise the type system will not allow one to construct a map. This
is the reason `Chain2` does not have `AsMap` method. Doing so would impose the
constraint to both `K` and `V` and limiting the usability of a `Chain2`.

```
invalid map key type K (missing comparable constraint)
```

{{ "ExampleAsMap" | example }}

All operations above can be chained.

{{ "ExampleChain2" | example }}

# Ideas

Some crazy and not so crazy ideas to expolse

## break the chain

One of the coolest (keep in mind it was a late night one) ideas may be breaking
the chain. The prototype exists in ideas_test.go, just not sure if it _is_
actually a good idea. It is definitely doable and possible in Go.

{{ "Example_break_da_chain" | example }}


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
