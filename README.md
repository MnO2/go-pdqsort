# go-pdqsort

[![Build Status](https://travis-ci.com/MnO2/go-pdqsort.svg?branch=master)](https://travis-ci.com/MnO2/go-pdqsort)
[![codecov](https://codecov.io/gh/MnO2/go-pdqsort/branch/master/graph/badge.svg)](https://codecov.io/gh/MnO2/go-pdqsort)

Pattern-defeating sort is a new variants of hybrid sort discovered in 2016, the
first proof of concept was in [C++](https://github.com/orlp/pdqsort) and later
ported to [rust](https://github.com/stjepang/pdqsort). It is now part of rust
standard libarary for unstable sort.

This library is to implement pattern-defeating sort in pure go. 

## Benchmark

Sort 1000 of Integers

* pdqsort

```
BenchmarkPDQSortInt1K-4   	   20000	     80874 ns/op	      32 B/op	       1 allocs/op
```

* sort from standard library

```
BenchmarkStdSortInt1K-4   	   20000	     69828 ns/op	      32 B/op	       1 allocs/op
```

The sorting algorithm in the standard library is also in hybrid approach. It
would go by shell sort on small input, and partition by quicksort up to a
certain recursive depth, then change to heapsort for better worst case scenario.

From the result the pattern defeating sort is much slower than the sort from
standrad library. It is expected since go's memory layout is mainly heap managed
and it is very different from what's in C++ and Rust. Therefore the key
invention from pdqsort to reduce the cpu branch prediction is not impactful on
speed, and slowing down due to more memory access. 
