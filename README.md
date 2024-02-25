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

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
res := it.NewMapable[string, int](it.From(n)).
	Map(func(s string) int { return len(s) }).
	Index().
	Filter2(func(index int, _ int) bool { return index <= 1 }).
	Values().
	Slice()
fmt.Println(res)
// Output: [2 3]
```

# Examples

## Filtering

Everything is available as a plain function. Some helpers are exposed via
`Chain` and `Mapable` structs methods allowing one to chain different
operations together. However this is more or less a syntax sugar on top regular
functions and explicit passing of variables.

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

The regular Go code uses a following patter for iterating through a slice

```go
for index, value := range slice {}
```

`it` does implement the `Index`/`IndexFrom` functions, which allows exactly that.

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

The way how this is implemented is it changes `iter.Seq[T]` into
`iter.Seq2[int, T]`. That has some consequences for a chain as the `Chain2`
struct is returned and `Filter2` method working on a pair must be called and
later `Values()` dropping the index from the sequence.

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

Map transforms one type to another. This was challenging to support in method
chaininig as Go type system don't allow to specify types of struct methods.
However it works for a common case mapping two types.

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

Supporting more than two types will lead to very messy API, however it is Go.
Old plain functions are always theanswer.

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

## Reduce

Reduce is a common functional operation, except it returns a single value. It
allows one to implement operation len.

```go
m := []int{1, 2, 3, 4, 5, 6, 7}
count := it.NewChain(it.From(m)).
	Reduce(func(a, _ int) int { return a + 1 }, 0)
fmt.Println(count)
// Output: 7
```

## Sort

All other operations can work on a single item at the time. Not sort - it first
pull all items to slice, sort it and then push the values to the iterator.

It accepts `type SortFunc[T any] func([]T)`, so caller can specify exactly
_how_ the sequence is going to be sorted. For example use a
`slices.SortStableFunc` to get a stable sort.

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
s0 := it.From(n)
s1 := it.Sort(s0, func(slice []string) { slices.SortFunc(slice, strings.Compare) })
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

## Map with errors

Sometimes there is no 1:1 transformation between `T` and `V` and mapping can
fail. For this reason `it` as a very generic mapping from `it.Seq[T]` into
`it.Seq2[K, V]`

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

As returning an error is so common operation in Go, there is a specialized function `MapError`

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

Which can be used inside a chain as well

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

## iter.Seq2[K, V]

Most operations does have the alternative working on `iter.Seq2[K, V]`. In
order to distinguish the names, all functions and chains has a suffix `2`, so
it is clear if method works with a single value or a pair.

Filtering - as the range order of iter.Seq2 is random in Go, the sequence must
be sorted first.

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

Sometimes one needs to pick only one of `K`, `V`, so `Keys` and `Values` does not have `2` suffix.

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

And one can get values back as a map - note the `K` must be `comparable`
otherwise type system does not allow one to construct a map.

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

## Chain2

All operations above can get chained.

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
