package runetrie

import (
	"errors"
	"unicode"
)

// ErrConflictEntry is returned when a new entry conflicts with an existing one.
// This can happen only in case insensitive mode.
// For example, if you add "foo" and then try to add "FOO", it will return this error.
var ErrConflictEntry = errors.New("conflict entry")

// Trie is a prefix tree (trie) .
// It is case sensitive by default.
type Trie[T ~string] struct {
	i bool // case insensitive
	m map[rune]*Trie[T]
	l struct {
		max, min int
	}
	s T
}

// Must is a helper function to create a new Trie and panic if an error occurs.
// It is useful for initializing a Trie in a single line.
// It should be used with caution, as it will panic on error.
// Use it only when you are sure that the error will not occur.
func Must[T ~string](trie *Trie[T], err error) *Trie[T] {
	if err != nil {
		panic(err)
	}
	return trie
}

// NewTrie creates a new case sensitive Trie with the given strings.
// It calls Add method to add the strings to the Trie internally.
func NewTrie[T ~string](ss ...T) *Trie[T] {
	trie := &Trie[T]{}
	_ = trie.Add(ss...)
	return trie
}

// NewCaseInsensitiveTrie creates a new case insensitive Trie with the given strings.
// It is useful for searching strings in a case insensitive manner.
// It calls Add method to add the strings to the Trie internally.
func NewCaseInsensitiveTrie[T ~string](ss ...T) (*Trie[T], error) {
	trie := &Trie[T]{i: true}
	if err := trie.Add(ss...); err != nil {
		return nil, err
	}
	return trie, nil
}

// Add adds a new string to the Trie.
// Only in case insensitive mode, it will check for conflicts.
// If a conflict is found, it returns ErrConflictEntry.
// If the string is already present, it does nothing and returns nil.
// If the string is not present, it adds it to the Trie and returns nil.
func (t *Trie[T]) Add(ss ...T) error {
	for _, s := range ss {
		if t.l.min == 0 || t.l.min > len(s) {
			t.l.min = len(s)
		}
		if t.l.max < len(s) {
			t.l.max = len(s)
		}

		tree := t
		for i, c := range s {
			if tree.m == nil {
				tree.m = map[rune]*Trie[T]{}
			}
			if leaf, ok := tree.m[c]; ok {
				if leaf.l.min == 0 || leaf.l.min > len(s[i+1:]) {
					leaf.l.min = len(s[i+1:])
				}
				if leaf.l.max < len(s[i+1:]) {
					leaf.l.max = len(s[i+1:])
				}
				tree = leaf
			} else {
				leaf := &Trie[T]{i: true}
				leaf.l.min = len(s[i+1:])
				leaf.l.max = len(s[i+1:])
				if t.i {
					tree.m[c] = leaf
					if unicode.IsLower(c) {
						tree.m[unicode.ToUpper(c)] = leaf
					} else if unicode.IsUpper(c) {
						tree.m[unicode.ToLower(c)] = leaf
					}
				} else {
					tree.m[c] = leaf
				}
				tree = leaf
			}
		}
		if tree.s != "" && tree.s != s {
			return ErrConflictEntry
		}
		tree.s = s
	}
	return nil
}

// MatchAny checks if any of the strings in the Trie match the given string.
// It returns true if there is a match, false otherwise.
func (t *Trie[T]) MatchAny(s T) bool {
	if t.m == nil {
		return false
	}

	tree := t
	for i, c := range s {
		if len(s[i:]) < tree.l.min || tree.l.max < len(s[i:]) {
			return false
		}

		if leaf, ok := tree.m[c]; ok {
			tree = leaf
			if tree.m == nil {
				break
			}
		} else {
			return false
		}
	}

	return tree.s != ""
}

// MatchAnyPrefixOf checks if any of the strings in the Trie match any prefix of the given string.
// It returns true if there is a match, false otherwise.
func (t *Trie[T]) MatchAnyPrefixOf(s T) bool {
	if len(s) < t.l.min {
		return false
	}
	if len(s) > t.l.max {
		s = s[:t.l.max]
	}

	tree := t
	for _, c := range s {
		if leaf, ok := tree.m[c]; ok {
			if leaf.s != "" {
				return true
			}
			tree = leaf
		} else {
			break
		}
	}
	return tree.s != ""
}

// MatchPrefixOf checks if the given string's prefix matches any of the strings in the Trie.
// It returns the shortest matched string and true if there is a match, or an empty string and false otherwise.
func (t *Trie[T]) MatchPrefixOf(s T) (T, bool) {
	if len(s) < t.l.min {
		var zero T
		return zero, false
	}
	if len(s) > t.l.max {
		s = s[:t.l.max]
	}

	tree := t
	for _, c := range s {
		if leaf, ok := tree.m[c]; ok {
			if leaf.s != "" {
				return leaf.s, true
			}
			tree = leaf
		} else {
			break
		}
	}

	var zero T
	return zero, false
}

// MatchPrefixOf checks if the given string's prefix matches any of the strings in the Trie.
// It returns the longest matched string and true if there is a match, or an empty string and false otherwise.
func (t *Trie[T]) LongestMatchPrefixOf(s T) (T, bool) {
	if len(s) < t.l.min {
		var zero T
		return zero, false
	}
	if len(s) > t.l.max {
		s = s[:t.l.max]
	}

	var result T
	matched := false
	tree := t
	for _, c := range s {
		if leaf, ok := tree.m[c]; ok {
			if leaf.s != "" {
				result = leaf.s
				matched = true
			}

			tree = leaf
			if tree.m == nil {
				break
			}
		} else {
			break
		}
	}
	return result, matched
}
