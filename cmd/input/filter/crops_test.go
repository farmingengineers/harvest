package filter

import (
	"slices"
	"testing"
)

func TestSplitCrop(t *testing.T) {
	examples := []struct {
		crop          string
		expectedParts []string
	}{
		{"", nil},
		{"word", []string{"word"}},
		{"word1 word2", []string{"word1", "word2"}},
		{"word1,word2", []string{"word1", "word2"}},
		{"word1, word2 word3", []string{"word1", "word2", "word3"}},
		{"squash, jack be little, each", []string{"squash", "jack", "be", "little", "each"}},
	}

	for _, ex := range examples {
		res := splitCrop(ex.crop)
		if !slices.Equal(ex.expectedParts, res) {
			t.Errorf("%q: expected %#v, got %#v", ex.crop, ex.expectedParts, res)
		}
	}
}

func BenchmarkSplitCrop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		splitCrop("squash, jack be little, each")
	}
}
