package prompts

// SecretFunc definition of the function for requesting a secret from the prompt.
type SecretFunc func(suppliedSecret string, prePromptNote string, promptLabel string, allowEmpty bool, confirmPrompt bool) (string, error)
