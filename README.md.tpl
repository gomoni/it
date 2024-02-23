# it: experimental iterator utilities for Go

# The problem

Go lacks ergonomic composable idiomatic iterator utility library. `it` builds
on top of a (rangefunc experiment)[https://go.dev/wiki/RangefuncExperiment] for
go 1.22. Design goals are

 * be idiomatic: so it provides all functionality as pure functions and use
   experimental `iter` package under the hood
 * be type safe: is uses generics everywhere
 * provide composable API when fits
 * support all Go's builtin types - slices/maps/channels

Non goals

> Supporting every possible iterator utility from
> lo/gubrak/Rust/Haskell/Scala/Python/whatever. Especially if those can be
> easily implemented via provided primitives.

IOW do not overwhelm users by all the functions it is possible to implement. Focus on
a real cases, which can't be easily built using primitives.

# Usage

> Don't forget to install go 1.22 and `export GOEXPERIMENT=rangefunc`

`it` provides methods, which can be chained together. Or a plain functions,
which covers more use cases. However method chains makes a better sales pitch

This example maps a slice to strings to int, add an index so only first two
items are returned and convert the code back to slice of ints.

{{ "ExampleMapable_Index" | example }}

# Examples

## Filtering

Everything is available as a plain function. Some helpers are exposed via
`Chain` and `Mapable` structs methods allowing one to chain different
operations together. However this is more or less a syntax sugar on top regular
functions and explicit passing of variables.

{{ "ExampleChain_Filter" | example }}

{{ "ExampleFilter" | example }}

## Indexing

The regular Go code uses a following patter for iterating through a slice

```go
for index, value := range slice {}
```

`it` does implement the `Index`/`IndexFrom` functions, which allows exactly that.

{{ "ExampleIndex" | example }}

The way how this is implemented is it changes `iter.Seq[T]` into
`iter.Seq2[int, T]`. That has some consequences for a chain as the `Chain2`
struct is returned and `Filter2` method working on a pair must be called and
later `Values()` dropping the index from the sequence.

{{ "ExampleChain_Index" | example }}

## Map

Map transforms one type to another. This was challenging to support in method
chaininig as Go type system don't allow to specify types of struct methods.
However it works for a common case mapping two types.

{{ "ExampleMapable_Map_second" | example }}

Supporting more than two types will lead to very messy API, however it is Go.
Old plain functions are always theanswer.

{{ "ExampleMap" | example }}

## Reduce

Reduce is a common functional operation, except it returns a single value. It
allows one to implement operation len.

{{ "ExampleChain_Reduce" | example }}

## Sort

All other operations can work on a single item at the time. Not sort - it first
pull all items to slice, sort it and then push the values to the iterator.

It accepts `type SortFunc[T any] func([]T)`, so caller can specify exactly
_how_ the sequence is going to be sorted. For example use a
`slices.SortStableFunc` to get a stable sort.

{{ "ExampleSort" | example }}

## Reverse

Simply iterate backward - it must collect the slice first.

{{ "ExampleReverse" | example }}

## Map with errors

Sometimes there is no 1:1 transformation between `T` and `V` and mapping can
fail. For this reason `it` as a very generic mapping from `it.Seq[T]` into
`it.Seq2[K, V]`

{{ "ExampleMapSeq2" | example }}

However returning an error is so common operation in go, there's simpler wrapper `MapError`

{{ "ExampleMapError" | example }}

> The chain support is not implemented in this case - will need one to rething it2.go struct(s)

## iter.Seq2[K, V]

Most operations does have the alternative working on `iter.Seq2[K, V]`. In
order to distinguish the names, all functions and chains has a suffix `2`, so
it is clear if method works with a single value or a pair.

Filtering - as the range order of iter.Seq2 is random in Go, the sequence must
be sorted first.

{{ "ExampleFilter2" | example }}

{{ "ExampleMap2" | example }}

Sometimes one needs to pick only one of `K`, `V`, so `Keys` and `Values` does not have `2` suffix.

{{ "ExampleKeys" | example }}

{{ "ExampleValues" | example }}

And one can get values back as a map - note the `K` must be `comparable`
otherwise type system does not allow one to construct a map.

> invalid map key type K (missing comparable constraint)

{{ "ExampleAsMap" | example }}

## Chain2

All operations above can get chained. The only limitation is `K` must be `comparable`.

{{ "ExampleChain2" | example }}

# WIP

## chain support for MapError

The biggest problem is Chain2 is typed as [K comparable, V any] and map is going to return `[V any, error]`

1. drop MapSeq2 - it's too complicated
2. reintroduce the WithError helper, which will return `[T, error]`
3. something else?

## break the chain

One of the coolest (keep in mind it was a late night one) ideas may be

```go
    n := []string{...}
    chain := it.NewChain(n).Filter().Something().Seq()
    // break the chain implementing a part as a for range loop
    for s, err := range chain {
        n, err := strconv.Atoi(s)
        if err != nil {
            break
        }
        // do the magic here and resume the chain
        magicYield(n)
    }

    Map(magicSeq, foo).Filter(bar).Slice()
```

Implemented in break_da_chain example test in ideas_test.go.

## Others

 * which operations `it` should have?
 * unit tests for operations on a (hash)map
 * what about operations which drains the input sequence? eg Keys?
 * naming and code organization - ie keep everything in a single `it` package and adopt `Seq` vs `Seq2` naming of a stdlib or move seq2 implementation into own package
 * make it context aware?????


# Other libraries

## https://pkg.go.dev/github.com/samber/lo

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


## https://github.com/novalagung/gubrak

Less popular, provides nicer API than lo. Implements own iterator type, so
methods can be arbitrary chained.

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
