# it: additional iterator utilities for Go

# The problem

Go 1.23 got [iterators](https://tip.golang.org/doc/go1.23#iterators) in a
standard library. `Filter/Map` operations were not included, so this library
adds them.

Additionally Go upstream prefers plain functions over methods, which makes a
chaining of iterator methods impossible. `it` provides a simple helper struct
`Chain` enabling operation chains like a `Filter/Map`.

# Usage

> Don't forget to install go 1.23rc2

```sh
$ go1.23rc2 get go@1.23rc2 toolchain@1.23rc2
$ export GOTOOLCHAIN=go1.23rc2
```

## Filter the slice

```go
n := []string{"aa", "aaa", "aaaaaaa", "a"}
s0 := slices.Values(n)
s1 := islices.Filter(s0, func(s string) bool { return len(s) >= 2 })
slice := slices.Collect(s1)
fmt.Println(slice)
// Output: [aa aaa aaaaaaa]
```

## Map the slice

```go
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	// maps string->int->float32
	s0 := slices.Values(n)
	s1 := islices.Map(s0, func(s string) int { return len(s) })
	s2 := islices.Map(s1, func(i int) float32 { return float32(i) })
	s3 := islices.Map(s2, func(f float32) string { return strconv.FormatFloat(float64(f), 'E', 4, 32) })
	slice := slices.Collect(s3)
	fmt.Println(slice)
	// Output: [2.0000E+00 3.0000E+00 7.0000E+00 1.0000E+00]
```

## Filter the map

```go
	m := map[string]int{
		"bambino": 1,
		"junior":  11,
		"senior":  22,
	}
	s0 := maps.All(m)
	s1 := imaps.Filter(s0, func(_ string, v int) bool { return v >= 18 })
	result := maps.Collect(s1)
	fmt.Println(result)
	// Output: map[senior:22]
```

## Map the map

```go
	m := map[string]int{
		"bambino": 1,
		"junior":  11,
		"senior":  22,
	}
	s0 := maps.All(m)
	s1 := imaps.Map(s0, func(s string, _ int) (string, int) { return s, len(s) })
	result := maps.Collect(s1)
	fmt.Println(result)
	// Output: map[bambino:7 junior:6 senior:6]
```

## Chainable filter

`Chain` provides API similar to other languages, where methods can be chained together

```go
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	ch := it.NewChain(slices.Values(n))
	slice := ch.
		Filter(func(s string) bool { return len(s) >= 2 }).
		Filter(func(s string) bool { return len(s) <= 4 }).
		Collect()
	fmt.Println(slice)
	// Output: [aa aaa]
```

## Chainable filter/Map

`Mappable` allows one to map from one type to another (and back) in a single chain

```go
	n := []string{"aa", "aaa", "aaaaaaa", "a"}
	ch := it.NewMapable[string, int](slices.Values(n))
	slice := ch.
		Filter(func(s string) bool { return len(s) >= 2 }).
		Map(func(s string) int { return len(s) }).
		Collect()
	fmt.Println(slice)
	// Output: [2 3 7]
```

# Performance

`slices.All` and `slices.Values` have a 50% performance impact. The combination
of a `slices.Values` and `islices.Filter` have much bigger impact thought.

```txt
go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/gomoni/it
cpu: AMD Ryzen 7 PRO 5850U with Radeon Graphics
BenchmarkRange-16                           4653            241190 ns/op
BenchmarkRangeAll-16                        2008            569001 ns/op
BenchmarkRangeValues-16                     2076            561859 ns/op
BenchmarkRangeEven-16                       4830            239241 ns/op
BenchmarkRangeValuesFilterEven-16            304           3940996 ns/op
PASS
ok      github.com/gomoni/it    6.362s
```

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
