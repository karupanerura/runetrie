package runetrie_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
	"github.com/karupanerura/runetrie"
)

func TestMust(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("must panic")
		} else if err != runetrie.ErrConflictEntry {
			t.Errorf("unexpected panic: %v", err)
		}
	}()
	runetrie.Must[string](nil, runetrie.ErrConflictEntry)
}

func TestNewCaseInsensitiveTrie_Conflict(t *testing.T) {
	tr, err := runetrie.NewCaseInsensitiveTrie("aA", "aa")
	if err == nil {
		t.Fatal("must be error")
	}
	if tr != nil {
		t.Errorf("must be omit trie: %+v", tr)
	}
	if !errors.Is(err, runetrie.ErrConflictEntry) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAdd(t *testing.T) {
	tr := runetrie.Must(runetrie.NewCaseInsensitiveTrie[string]())
	err := tr.Add("Aa", "aA")
	if err == nil {
		t.Fatal("must be error")
	}
	if !errors.Is(err, runetrie.ErrConflictEntry) {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_Trie_MatchAnyPrefixOf(t *testing.T) {
	tests := []struct {
		name   string
		set    []string
		target string
		want   bool
	}{
		{
			name:   "Empty",
			set:    []string{},
			target: "",
			want:   false,
		},
		{
			name:   "EmptyIsNotAnyStrings",
			set:    []string{},
			target: "foo",
			want:   false,
		},
		{
			name:   "EmptyIsNotMatchAnyStrings",
			set:    []string{"foo"},
			target: "",
			want:   false,
		},
		{
			name:   "ExactlyMatchA",
			set:    []string{"A"},
			target: "A",
			want:   true,
		},
		{
			name:   "ExactlyMatchAAA",
			set:    []string{"A", "AA", "AAA"},
			target: "AAA",
			want:   true,
		},
		{
			name:   "PrefixMatchAAC",
			set:    []string{"A", "AA", "AAA"},
			target: "AAC",
			want:   true,
		},
		{
			name:   "DoNotContainNotCompletedTree",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "ABC",
			want:   false,
		},
		{
			name:   "CaseSensitive",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "abcabc",
			want:   false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := runetrie.NewTrie(tt.set...)
			got := tr.MatchAnyPrefixOf(tt.target)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Trie.MatchAnyPrefixOf() = %v, want %v.\n%s", got, tt.want, diff)
				t.Log(pp.Sprint(tr))
			}
		})
	}
}

func Test_Trie_MatchPrefixOf(t *testing.T) {
	type ret struct {
		Result  string
		Matched bool
	}
	tests := []struct {
		name   string
		set    []string
		target string
		want   ret
	}{
		{
			name:   "Empty",
			set:    []string{},
			target: "",
			want:   ret{"", false},
		},
		{
			name:   "EmptyIsNotAnyStrings",
			set:    []string{},
			target: "foo",
			want:   ret{"", false},
		},
		{
			name:   "EmptyIsNotMatchAnyStrings",
			set:    []string{"foo"},
			target: "",
			want:   ret{"", false},
		},
		{
			name:   "ExactlyMatchA",
			set:    []string{"A"},
			target: "A",
			want:   ret{"A", true},
		},
		{
			name:   "ExactlyMatchAAA",
			set:    []string{"A", "AA", "AAA"},
			target: "AAA",
			want:   ret{"A", true},
		},
		{
			name:   "TooShort",
			set:    []string{"AA", "AAA"},
			target: "A",
			want:   ret{"", false},
		},
		{
			name:   "Mismatch",
			set:    []string{"AA", "AAA"},
			target: "ABC",
			want:   ret{"", false},
		},
		{
			name:   "PrefixMatchAAC",
			set:    []string{"A", "AA", "AAA"},
			target: "AAC",
			want:   ret{"A", true},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := runetrie.NewTrie(tt.set...)
			result, matched := tr.MatchPrefixOf(tt.target)
			if diff := cmp.Diff(tt.want, ret{result, matched}); diff != "" {
				t.Errorf("Trie.MatchPrefixOf() = (%v %v), want %v.\n%s", result, matched, tt.want, diff)
				t.Log(pp.Sprint(tr))
			}
		})
	}
}

func Test_Trie_LongestMatchPrefixOf(t *testing.T) {
	type ret struct {
		Result  string
		Matched bool
	}
	tests := []struct {
		name   string
		set    []string
		target string
		want   ret
	}{
		{
			name:   "Empty",
			set:    []string{},
			target: "",
			want:   ret{"", false},
		},
		{
			name:   "EmptyIsNotAnyStrings",
			set:    []string{},
			target: "foo",
			want:   ret{"", false},
		},
		{
			name:   "EmptyIsNotMatchAnyStrings",
			set:    []string{"foo"},
			target: "",
			want:   ret{"", false},
		},
		{
			name:   "ExactlyMatchA",
			set:    []string{"A"},
			target: "A",
			want:   ret{"A", true},
		},
		{
			name:   "ExactlyMatchAAA",
			set:    []string{"A", "AA", "AAA"},
			target: "AAA",
			want:   ret{"AAA", true},
		},
		{
			name:   "PrefixMatchAAC",
			set:    []string{"A", "AA", "AAA"},
			target: "AAC",
			want:   ret{"AA", true},
		},
		{
			name:   "CaseInsensitive",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "abcabc",
			want:   ret{"", false},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := runetrie.NewTrie(tt.set...)
			result, matched := tr.LongestMatchPrefixOf(tt.target)
			if diff := cmp.Diff(tt.want, ret{result, matched}); diff != "" {
				t.Errorf("Trie.LongestMatchPrefixOf() = (%v, %v), want %v.\n%s", result, matched, tt.want, diff)
				t.Log(pp.Sprint(tr))
			}
		})
	}
}

func Test_Trie_MatchAny(t *testing.T) {
	tests := []struct {
		name   string
		set    []string
		target string
		want   bool
	}{
		{
			name:   "Empty",
			set:    []string{},
			target: "",
			want:   false,
		},
		{
			name:   "EmptyIsNotAnyStrings",
			set:    []string{},
			target: "foo",
			want:   false,
		},
		{
			name:   "EmptyIsNotMatchAnyStrings",
			set:    []string{"foo"},
			target: "",
			want:   false,
		},
		{
			name:   "ExactlyMatchA",
			set:    []string{"A"},
			target: "A",
			want:   true,
		},
		{
			name:   "ExactlyMatchAAA",
			set:    []string{"A", "AA", "AAA"},
			target: "AAA",
			want:   true,
		},
		{
			name:   "PrefixMatchAAC",
			set:    []string{"A", "AA", "AAA"},
			target: "AAC",
			want:   false,
		},
		{
			name:   "DoNotContainNotCompletedTree",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "ABC",
			want:   false,
		},
		{
			name:   "CaseSensitive",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "abca",
			want:   false,
		},
		{
			name:   "Symbols",
			set:    []string{">", ">=", "<", "<="},
			target: ">",
			want:   true,
		},
		{
			name:   "EarlyExit",
			set:    []string{"ABC", "ABCDE", "ABCDEF"},
			target: "ABCD",
			want:   false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := runetrie.NewTrie(tt.set...)
			if got := tr.MatchAny(tt.target); got != tt.want {
				t.Errorf("Trie.MatchAny() = %v, want %v", got, tt.want)
				t.Log(pp.Sprint(tr))
			}
		})
	}
}

func Test_CaseInsensitiveTrie_MatchAnyPrefixOf(t *testing.T) {
	tests := []struct {
		name   string
		set    []string
		target string
		want   bool
	}{
		{
			name:   "Empty",
			set:    []string{},
			target: "",
			want:   false,
		},
		{
			name:   "EmptyIsNotAnyStrings",
			set:    []string{},
			target: "foo",
			want:   false,
		},
		{
			name:   "EmptyIsNotMatchAnyStrings",
			set:    []string{"foo"},
			target: "",
			want:   false,
		},
		{
			name:   "ExactlyMatchA",
			set:    []string{"A"},
			target: "A",
			want:   true,
		},
		{
			name:   "ExactlyMatchAAA",
			set:    []string{"A", "AA", "AAA"},
			target: "AAA",
			want:   true,
		},
		{
			name:   "PrefixMatchAAC",
			set:    []string{"A", "AA", "AAA"},
			target: "AAC",
			want:   true,
		},
		{
			name:   "DoNotContainNotCompletedTree",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "ABC",
			want:   false,
		},
		{
			name:   "CaseInsensitive",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "abcabc",
			want:   true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := runetrie.Must(runetrie.NewCaseInsensitiveTrie(tt.set...))
			got := tr.MatchAnyPrefixOf(tt.target)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Trie.MatchAnyPrefixOf() = %v, want %v.\n%s", got, tt.want, diff)
				t.Log(pp.Sprint(tr))
			}
		})
	}
}

func Test_CaseInsensitiveTrie_MatchPrefixOf(t *testing.T) {
	type ret struct {
		Result  string
		Matched bool
	}
	tests := []struct {
		name   string
		set    []string
		target string
		want   ret
	}{
		{
			name:   "Empty",
			set:    []string{},
			target: "",
			want:   ret{"", false},
		},
		{
			name:   "EmptyIsNotAnyStrings",
			set:    []string{},
			target: "foo",
			want:   ret{"", false},
		},
		{
			name:   "EmptyIsNotMatchAnyStrings",
			set:    []string{"foo"},
			target: "",
			want:   ret{"", false},
		},
		{
			name:   "ExactlyMatchA",
			set:    []string{"A"},
			target: "A",
			want:   ret{"A", true},
		},
		{
			name:   "ExactlyMatchAAA",
			set:    []string{"A", "AA", "AAA"},
			target: "AAA",
			want:   ret{"A", true},
		},
		{
			name:   "PrefixMatchAAC",
			set:    []string{"A", "AA", "AAA"},
			target: "AAC",
			want:   ret{"A", true},
		},
		{
			name:   "CaseInsensitive",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "abcabc",
			want:   ret{"ABCA", true},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := runetrie.Must(runetrie.NewCaseInsensitiveTrie(tt.set...))
			result, matched := tr.MatchPrefixOf(tt.target)
			if diff := cmp.Diff(tt.want, ret{result, matched}); diff != "" {
				t.Errorf("Trie.MatchPrefixOf() = (%v %v), want %v.\n%s", result, matched, tt.want, diff)
				t.Log(pp.Sprint(tr))
			}
		})
	}
}

func Test_CaseInsensitiveTrie_LongestMatchPrefixOf(t *testing.T) {
	type ret struct {
		Result  string
		Matched bool
	}
	tests := []struct {
		name   string
		set    []string
		target string
		want   ret
	}{
		{
			name:   "Empty",
			set:    []string{},
			target: "",
			want:   ret{"", false},
		},
		{
			name:   "EmptyIsNotAnyStrings",
			set:    []string{},
			target: "foo",
			want:   ret{"", false},
		},
		{
			name:   "EmptyIsNotMatchAnyStrings",
			set:    []string{"foo"},
			target: "",
			want:   ret{"", false},
		},
		{
			name:   "ExactlyMatchA",
			set:    []string{"A"},
			target: "A",
			want:   ret{"A", true},
		},
		{
			name:   "ExactlyMatchAAA",
			set:    []string{"A", "AA", "AAA"},
			target: "AAA",
			want:   ret{"AAA", true},
		},
		{
			name:   "PrefixMatchAAC",
			set:    []string{"A", "AA", "AAA"},
			target: "AAC",
			want:   ret{"AA", true},
		},
		{
			name:   "CaseInsensitive",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "abcabc",
			want:   ret{"ABCA", true},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := runetrie.Must(runetrie.NewCaseInsensitiveTrie(tt.set...))
			result, matched := tr.LongestMatchPrefixOf(tt.target)
			if diff := cmp.Diff(tt.want, ret{result, matched}); diff != "" {
				t.Errorf("Trie.LongestMatchPrefixOf() = (%v %v), want %v.\n%s", result, matched, tt.want, diff)
				t.Log(pp.Sprint(tr))
			}
		})
	}
}

func Test_CaseInsensitiveTrie_MatchAny(t *testing.T) {
	tests := []struct {
		name   string
		set    []string
		target string
		want   bool
	}{
		{
			name:   "Empty",
			set:    []string{},
			target: "",
			want:   false,
		},
		{
			name:   "EmptyIsNotAnyStrings",
			set:    []string{},
			target: "foo",
			want:   false,
		},
		{
			name:   "EmptyIsNotMatchAnyStrings",
			set:    []string{"foo"},
			target: "",
			want:   false,
		},
		{
			name:   "ExactlyMatchA",
			set:    []string{"A"},
			target: "A",
			want:   true,
		},
		{
			name:   "ExactlyMatchAAA",
			set:    []string{"A", "AA", "AAA"},
			target: "AAA",
			want:   true,
		},
		{
			name:   "PrefixMatchAAC",
			set:    []string{"A", "AA", "AAA"},
			target: "AAC",
			want:   false,
		},
		{
			name:   "DoNotContainNotCompletedTree",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "ABC",
			want:   false,
		},
		{
			name:   "CaseInsensitive",
			set:    []string{"AAAA", "ABAA", "ACAA", "ABCA"},
			target: "abca",
			want:   true,
		},
		{
			name:   "Symbols",
			set:    []string{">", ">=", "<", "<="},
			target: ">",
			want:   true,
		},
		{
			name:   "EarlyExit",
			set:    []string{"ABC", "ABCDE", "ABCDEF"},
			target: "ABCD",
			want:   false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := runetrie.Must(runetrie.NewCaseInsensitiveTrie(tt.set...))
			if got := tr.MatchAny(tt.target); got != tt.want {
				t.Errorf("Trie.MatchAny() = %v, want %v", got, tt.want)
				t.Log(pp.Sprint(tr))
			}
		})
	}
}
