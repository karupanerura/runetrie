package runetrie_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/karupanerura/runetrie"
)

var (
	targetStrings []string
	searchStrings []string
)

func init() {
	for _, c := range "ABC" {
		r := []rune{c}
		targetStrings = append(targetStrings, string(r))
		for _, c := range "ABC" {
			r = append(r, c)
			targetStrings = append(targetStrings, string(r))
			for _, c := range "ABC" {
				r = append(r, c)
				targetStrings = append(targetStrings, string(r))
			}
		}
	}
	sort.Strings(targetStrings)

	searchStrings = make([]string, len(targetStrings))
	copy(searchStrings, targetStrings)
	for _, c := range "DEF" {
		r := []rune{c}
		targetStrings = append(targetStrings, string(r))
		for _, c := range "DEF" {
			r = append(r, c)
			targetStrings = append(targetStrings, string(r))
			for _, c := range "DEF" {
				r = append(r, c)
				targetStrings = append(targetStrings, string(r))
			}
		}
	}
	sort.Strings(searchStrings)
}

func BenchmarkStringsHasPrefixShortestMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s1 := range searchStrings {
			for _, s2 := range targetStrings {
				if strings.HasPrefix(s2, s1) {
					break
				}
			}
		}
	}
}

func BenchmarkStringsHasPrefixLongestMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s1 := range searchStrings {
			for i := len(targetStrings) - 1; i >= 0; i-- {
				if strings.HasPrefix(targetStrings[i], s1) {
					break
				}
			}
		}
	}
}

func BenchmarkStringsMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s1 := range searchStrings {
			for _, s2 := range targetStrings {
				if s1 == s2 {
					break
				}
			}
		}
	}
}

func BenchmarkTrieMatchAny(b *testing.B) {
	trie := runetrie.NewTrie(targetStrings...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, s := range searchStrings {
			trie.MatchAny(s)
		}
	}
}

func BenchmarkTrieMatchAnyPrefixOf(b *testing.B) {
	trie := runetrie.NewTrie(targetStrings...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, s := range searchStrings {
			trie.MatchAnyPrefixOf(s)
		}
	}
}

func BenchmarkTrieMatchPrefixOf(b *testing.B) {
	trie := runetrie.NewTrie(targetStrings...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, s := range searchStrings {
			trie.MatchPrefixOf(s)
		}
	}
}

func BenchmarkTrieLongestMatchPrefixOf(b *testing.B) {
	trie := runetrie.NewTrie(targetStrings...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, s := range searchStrings {
			trie.LongestMatchPrefixOf(s)
		}
	}
}
