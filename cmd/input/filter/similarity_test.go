package filter

import (
	"fmt"
	"testing"
)

func TestSimilarity(t *testing.T) {
	examples := []struct {
		input string
		more  string
		less  []string
	}{
		{
			input: "cu",
			more:  "cucumbers",
			less: []string{
				"cabbage",
				"carrots",
				"celeriac",
			},
		},
	}

	for _, ex := range examples {
		for _, less := range ex.less {
			t.Run(fmt.Sprintf("%s/%s/%s", ex.input, ex.more, less), func(t *testing.T) {
				a := jaroWinklerSimilarity(ex.input, ex.more)
				b := jaroWinklerSimilarity(ex.input, less)
				if a < b {
					t.Errorf("Expected %q (score=%f) to be more similar than %q (score=%f)",
						ex.more, a,
						less, b)
				}
			})
		}
	}
}
