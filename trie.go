package runetrie

import (
	"errors"
	"unicode"
)

var ErrConflictEntry = errors.New("conflict entry")

type Trie[T ~string] struct {
	i bool // case insensitive
	m map[rune]*Trie[T]
	l struct {
		max, min int
	}
	s T
}

func Must[T ~string](trie *Trie[T], err error) *Trie[T] {
	if err != nil {
		panic(err)
	}
	return trie
}

func NewTrie[T ~string](ss ...T) *Trie[T] {
	trie := &Trie[T]{}
	_ = trie.Add(ss...)
	return trie
}

func NewCaseInsensitiveTrie[T ~string](ss ...T) (*Trie[T], error) {
	trie := &Trie[T]{i: true}
	if err := trie.Add(ss...); err != nil {
		return nil, err
	}
	return trie, nil
}

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

func (t *Trie[T]) MatchAny(s T) bool {
	if t.m == nil {
		return false
	}
	if len(s) < t.l.min || t.l.max < len(s) {
		return false
	}

	tree := t
	for i, c := range s {
		if leaf, ok := tree.m[c]; ok {
			if len(s[i+1:]) < leaf.l.min || leaf.l.max < len(s[i+1:]) {
				break
			}

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
