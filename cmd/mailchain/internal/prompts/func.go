package prompts

type SecretFunc func(suppliedSecret string, prePromptNote string, promptLabel string, allowEmpty bool, confirmPrompt bool) (string, error)
