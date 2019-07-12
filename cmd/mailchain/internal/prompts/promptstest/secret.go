package promptstest

import "testing"

func MockRequiredSecret(t *testing.T, returnValue string, returnErr error) func(suppliedSecret, prePromptNote, promptLabel string, allowEmpty, confirmPrompt bool) (string, error) {
	return func(suppliedSecret, prePromptNote, promptLabel string, allowEmpty, confirmPrompt bool) (string, error) {
		return returnValue, returnErr
	}
}
