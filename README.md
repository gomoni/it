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

IOW do not overwhelm users by all the functions it can provide.

# Usage

`it` provides methods, which can be chained together. Or a plain functions,
which covers more use cases. However method chains makes a great sales pitch

```go
	n := []string{"aa", "aaa", "aaaaaaa", "a"}

	// as well as a map T -> V
	res := it.NewMappable[string, int](n).
		Map(func(s string) int { return len(s) }).
		Filter(func(i int) bool { return i >= 2 }).
		Slice()
	fmt.Println(res)
	// Output: [2 3 7]
```

Don't forget to install go 1.22 and `export GOEXPERIMENT=rangefunc`

# WIP

## How enumerable stuff or errors?

It turns out there are two equivalent solutions for both

## "upgrade" to `iter.Seq2`

As shown in `Example_idea_errors` - this is probably more idiomatic solution
for index numbers than errors as `iter.Seq2` is a direct equivalent of

```go
for index, value := range slice {}
```

## provide a wrapper struct

As shown in `Example_idea_enumerable` the whole problem can be solved by simply

1. wrapper struct `Indexed`
2. `Map(enumerable)` maps the sequence into sequence with indices

## Conclusion

Maybe the best solution is to provide both

1. an idiomatic way how to "upgrade" the iter.Seq into iter.Seq2
2. provide a default wrappers for common cases like `Indexed`, `Fallible`

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


* Chunk     -> Map and func(T) []T
* Compact   -> Filter
* Concat
* Count     -> implement Reduce
* Difference
* DifferenceMany
* Drop      -> Drop can be implemented via enumerable Filter
* DropRight -> ditto
* Each      -> interesting in supporting early end
* EachRight 
* Fill 
* Filter 
* Find 
* FindIndex
* FindLast 
* FindLastIndex 
* First 
* ForEach 
* ForEachRight 
* FromPairs 
* GroupBy 
* Head 
* Includes 
* IndexOf 
* Initial 
* Intersection 
* Join 
* KeyBy 
* Last 
* LastIndexOf 
* Map 
* Now 
* Nth 
* OrderBy 
* Partition 
* Pull 
* PullAll 
* PullAt 
* RandomInt 
* RandomString 
* Reduce 
* Reject 
* Remove 
* Reverse 
* Sample 
* SampleSize 
* Shuffle 
* Size 
* SortBy 
* Tail 
* Take 
* TakeRight 
* Union 
* Uniq
* Without

## https://docs.python.org/3/library/itertools.html
## https://docs.python.org/3/library/functools.html#functools.reduce
## https://doc.rust-lang.org/book/ch13-02-iterators.html
## https://hackage.haskell.org/package/base-4.19.0.0/docs/Data-List.html#g:2
