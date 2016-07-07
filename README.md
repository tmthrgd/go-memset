# go-memset

[![GoDoc](https://godoc.org/github.com/tmthrgd/go-memset?status.svg)](https://godoc.org/github.com/tmthrgd/go-memset)
[![Build Status](https://travis-ci.org/tmthrgd/go-memset.svg?branch=master)](https://travis-ci.org/tmthrgd/go-memset)

An efficient memset implementation for Golang.

In Golang the following loop is optimised with an assembly implementation in src/runtime/memclr_$GOARCH.s
(since [137880043](https://golang.org/cl/137880043)):
```
for i := range data {
	data[i] = 0
}
```
but the following loop is *not* optimised:
```
for i := range data {
	data[i] = 0xff
}
```
and neither is:
```
for i := range data {
	data[i] = value
}
```

go-memset provides a Memset function which uses an assembly implementation on x86-64 and can provide
performance equivalent to the optimised first loop.

## Download

```
go get github.com/tmthrgd/go-memset
```

## Benchmark

```
BenchmarkZero_32-8    	200000000	         6.38 ns/op	5013.17 MB/s
BenchmarkZero_128-8   	200000000	         6.67 ns/op	19198.09 MB/s
BenchmarkZero_1k-8    	50000000	        23.1 ns/op	44302.00 MB/s
BenchmarkZero_16k-8   	 5000000	       278 ns/op	58833.66 MB/s
BenchmarkZero_128k-8  	  500000	      3703 ns/op	35389.69 MB/s
BenchmarkZero_1M-8    	   50000	     39952 ns/op	26245.73 MB/s
BenchmarkZero_16M-8   	    2000	   1100177 ns/op	15249.55 MB/s
BenchmarkZero_128M-8  	     200	   8791286 ns/op	15267.13 MB/s
BenchmarkZero_512M-8  	      50	  35738899 ns/op	15022.03 MB/s
BenchmarkSet_32-8     	200000000	         6.33 ns/op	5058.06 MB/s
BenchmarkSet_128-8    	200000000	         6.56 ns/op	19522.67 MB/s
BenchmarkSet_1k-8     	50000000	        22.9 ns/op	44785.40 MB/s
BenchmarkSet_16k-8    	 5000000	       278 ns/op	58782.02 MB/s
BenchmarkSet_128k-8   	  500000	      3725 ns/op	35183.52 MB/s
BenchmarkSet_1M-8     	   30000	     40357 ns/op	25981.87 MB/s
BenchmarkSet_16M-8    	    2000	   1119679 ns/op	14983.94 MB/s
BenchmarkSet_128M-8   	     200	   8764132 ns/op	15314.43 MB/s
BenchmarkSet_512M-8   	      50	  35856688 ns/op	14972.69 MB/s

BenchmarkZeroGo_32-8  	300000000	         5.49 ns/op	5832.67 MB/s
BenchmarkZeroGo_128-8 	200000000	         6.49 ns/op	19736.93 MB/s
BenchmarkZeroGo_1k-8  	100000000	        21.1 ns/op	48610.78 MB/s
BenchmarkZeroGo_16k-8 	 5000000	       282 ns/op	57925.53 MB/s
BenchmarkZeroGo_128k-8	  500000	      3744 ns/op	35007.18 MB/s
BenchmarkZeroGo_1M-8  	   30000	     40439 ns/op	25929.63 MB/s
BenchmarkZeroGo_16M-8 	    1000	   1317314 ns/op	12735.92 MB/s
BenchmarkZeroGo_128M-8	     100	  10727737 ns/op	12511.28 MB/s
BenchmarkZeroGo_512M-8	      30	  43004291 ns/op	12484.12 MB/s
BenchmarkSetGo_32-8   	50000000	        28.5 ns/op	1123.90 MB/s
BenchmarkSetGo_128-8  	20000000	        83.5 ns/op	1532.35 MB/s
BenchmarkSetGo_1k-8   	 3000000	       555 ns/op	1842.64 MB/s
BenchmarkSetGo_16k-8  	  200000	      8651 ns/op	1893.79 MB/s
BenchmarkSetGo_128k-8 	   20000	     69114 ns/op	1896.44 MB/s
BenchmarkSetGo_1M-8   	    3000	    552932 ns/op	1896.39 MB/s
BenchmarkSetGo_16M-8  	     200	   8950601 ns/op	1874.42 MB/s
BenchmarkSetGo_128M-8 	      20	  71233792 ns/op	1884.19 MB/s
BenchmarkSetGo_512M-8 	       5	 285058677 ns/op	1883.37 MB/s
```

```
benchmark                Go ns/op      asm ns/op     delta
BenchmarkZero_32-8       5.49          6.38          +16.21%
BenchmarkZero_128-8      6.49          6.67          +2.77%
BenchmarkZero_1k-8       21.1          23.1          +9.48%
BenchmarkZero_16k-8      282           278           -1.42%
BenchmarkZero_128k-8     3744          3703          -1.10%
BenchmarkZero_1M-8       40439         39952         -1.20%
BenchmarkZero_16M-8      1317314       1100177       -16.48%
BenchmarkZero_128M-8     10727737      8791286       -18.05%
BenchmarkZero_512M-8     43004291      35738899      -16.89%
BenchmarkSet_32-8        28.5          6.33          -77.79%
BenchmarkSet_128-8       83.5          6.56          -92.14%
BenchmarkSet_1k-8        555           22.9          -95.87%
BenchmarkSet_16k-8       8651          278           -96.79%
BenchmarkSet_128k-8      69114         3725          -94.61%
BenchmarkSet_1M-8        552932        40357         -92.70%
BenchmarkSet_16M-8       8950601       1119679       -87.49%
BenchmarkSet_128M-8      71233792      8764132       -87.70%
BenchmarkSet_512M-8      285058677     35856688      -87.42%

benchmark                Go MB/s      asm MB/s     speedup
BenchmarkZero_32-8       5832.67      5013.17      0.86x
BenchmarkZero_128-8      19736.93     19198.09     0.97x
BenchmarkZero_1k-8       48610.78     44302.00     0.91x
BenchmarkZero_16k-8      57925.53     58833.66     1.02x
BenchmarkZero_128k-8     35007.18     35389.69     1.01x
BenchmarkZero_1M-8       25929.63     26245.73     1.01x
BenchmarkZero_16M-8      12735.92     15249.55     1.20x
BenchmarkZero_128M-8     12511.28     15267.13     1.22x
BenchmarkZero_512M-8     12484.12     15022.03     1.20x
BenchmarkSet_32-8        1123.90      5058.06      4.50x
BenchmarkSet_128-8       1532.35      19522.67     12.74x
BenchmarkSet_1k-8        1842.64      44785.40     24.31x
BenchmarkSet_16k-8       1893.79      58782.02     31.04x
BenchmarkSet_128k-8      1896.44      35183.52     18.55x
BenchmarkSet_1M-8        1896.39      25981.87     13.70x
BenchmarkSet_16M-8       1874.42      14983.94     7.99x
BenchmarkSet_128M-8      1884.19      15314.43     8.13x
BenchmarkSet_512M-8      1883.37      14972.69     7.95x

```

## License

Unless otherwise noted, the go-popcount source files are distributed under the Modified BSD License
found in the LICENSE file.
