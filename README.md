# runetrie

[![Go Reference](https://pkg.go.dev/badge/github.com/karupanerura/runetrie.svg)](https://pkg.go.dev/github.com/karupanerura/runetrie)

Package runetrie is a Go implementation of a trie (prefix tree) for efficient string matching. It uses Go's generics to work with any string-like type and performs better than standard string comparison methods when dealing with large string sets, as shown in the benchmarks below.

## Benchmark

I benchmarked this package `MatchAny`, `MatchAnyPrefixOf` `MatchPrefixOf` and `LongestMatchPrefixOf` as well as its performance against the standard `==` and `strings.HasPrefix`. The results are as follows.

```
goos: darwin
goarch: arm64
pkg: github.com/karupanerura/runetrie
BenchmarkStringsHasPrefixShortestMatch-12    	  721148	      1662 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringsHasPrefixLongestMatch-12     	  235120	      5099 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringsMatch-12                     	 1912580	       624.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkTrieMatchAny-12                     	  712273	      1701 ns/op	       0 B/op	       0 allocs/op
BenchmarkTrieMatchAnyPrefixOf-12             	 4555736	       263.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkTrieMatchPrefixOf-12                	 4533429	       263.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkTrieLongestMatchPrefixOf-12         	  733917	      1647 ns/op	       0 B/op	       0 allocs/op
```

## License

MIT License