package promptstest

import (
	"testing"
)

func MockRequiredInputWithDefault(t *testing.T, returnValue string, returnErr error) func(label string, defaultValue string) (string, error) {
	return func(label string, defaultValue string) (string, error) {
		return returnValue, returnErr
	}
}
