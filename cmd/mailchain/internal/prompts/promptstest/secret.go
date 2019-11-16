package promptstest

import "testing"

// MockRequiredSecret a mock for ensuring RequiredSecret is called in tests.
func MockRequiredSecret(t *testing.T, returnValue string, returnErr error) func(suppliedSecret, prePromptNote, promptLabel string, allowEmpty, confirmPrompt bool) (string, error) {
	return func(suppliedSecret, prePromptNote, promptLabel string, allowEmpty, confirmPrompt bool) (string, error) {
		return returnValue, returnErr
	}
}
