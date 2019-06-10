package promptstest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockSelectItemSkipable(t *testing.T, wantItems []string, returnSelected string, returnSkipped bool, returnErr error) func(label string, items []string, skipable bool) (selected string, skipped bool, err error) {
	return func(label string, items []string, skipable bool) (selected string, skipped bool, err error) {
		if !assert.EqualValues(t, wantItems, items) {
			t.Errorf("items = %v, wantItems %v", items, wantItems)
		}

		return returnSelected, returnSkipped && skipable, returnErr
	}
}
