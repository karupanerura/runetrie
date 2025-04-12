package runetrie_test

import (
	"fmt"

	"github.com/karupanerura/runetrie"
)

func ExampleNewTrie() {
	trie := runetrie.NewTrie("foo", "bar", "buz")
	fmt.Println(trie.MatchAny("foo"))
	// Output: true
}

func ExampleNewCaseInsensitiveTrie() {
	trie := runetrie.Must(runetrie.NewCaseInsensitiveTrie("FoO", "bAr", "buz"))
	fmt.Println(trie.LongestMatchPrefixOf("fOoo"))
	fmt.Println(trie.LongestMatchPrefixOf("BaR!"))
	// Output:
	// FoO true
	// bAr true
}

func ExampleTrie_Add() {
	trie := runetrie.NewTrie("foo", "bar", "buz")
	fmt.Println(trie.MatchAny("hoge"))
	trie.Add("hoge")
	fmt.Println(trie.MatchAny("hoge"))
	// Output:
	// false
	// true
}

func ExampleTrie_MatchAny() {
	trie := runetrie.NewTrie("foo", "bar", "buz")
	fmt.Println(trie.MatchAny("fo"))
	fmt.Println(trie.MatchAny("foo"))
	fmt.Println(trie.MatchAny("fooo"))
	// Output:
	// false
	// true
	// false
}

func ExampleTrie_MatchAnyPrefixOf() {
	trie := runetrie.NewTrie("foo", "bar", "buz")
	fmt.Println(trie.MatchAnyPrefixOf("fo"))
	fmt.Println(trie.MatchAnyPrefixOf("foo"))
	fmt.Println(trie.MatchAnyPrefixOf("fooo"))
	// Output:
	// false
	// true
	// true
}

func ExampleTrie_MatchPrefixOf() {
	trie := runetrie.NewTrie("application/json", "application/xml", "text/html", "text/plain")
	matched, ok := trie.MatchPrefixOf("application/json; charset=utf-8")
	if ok {
		fmt.Println(matched)
	} else {
		fmt.Println("No match found")
	}

	// Output:
	// application/json
}

func ExampleTrie_MatchPrefixOf_shortest() {
	trie := runetrie.NewTrie("AA", "AB", "AC", "AD", "AAA", "ABA", "ABC", "ABBB", "ACCC", "ADDD")
	matched, ok := trie.MatchPrefixOf("ABCD")
	if ok {
		fmt.Println(matched)
	} else {
		fmt.Println("No match found")
	}

	// Output:
	// AB
}

func ExampleTrie_LongestMatchPrefixOf() {
	trie := runetrie.NewTrie("application/json", "application/xml", "text/html", "text/plain")
	matched, ok := trie.LongestMatchPrefixOf("application/json; charset=utf-8")
	if ok {
		fmt.Println(matched)
	} else {
		fmt.Println("No match found")
	}

	// Output:
	// application/json
}

func ExampleTrie_LongestMatchPrefixOf_longest() {
	trie := runetrie.NewTrie("AA", "AB", "AC", "AD", "AAA", "ABA", "ABC", "ABBB", "ACCC", "ADDD")
	matched, ok := trie.LongestMatchPrefixOf("ABCD")
	if ok {
		fmt.Println(matched)
	} else {
		fmt.Println("No match found")
	}

	// Output:
	// ABC
}
