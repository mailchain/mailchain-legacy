package prompts

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

// Secret is extracted from the command, otherwise from prompt
func Secret(suppliedSecret, prePromptNote, promptLabel string, allowEmpty, confirmPrompt bool) (string, error) {
	if suppliedSecret != "" {
		return suppliedSecret, nil
	}
	if allowEmpty {
		return "", nil
	}
	fmt.Println(prePromptNote)
	return secretFromPrompt(promptLabel, confirmPrompt)
}

func secretFromPrompt(promptLabel string, confirmPrompt bool) (string, error) {
	prompt := promptui.Prompt{
		Label: promptLabel,
		Mask:  '*',
	}
	secret, err := prompt.Run()
	if err != nil {
		return "", errors.Errorf("failed read %q", promptLabel)
	}
	if confirmPrompt {
		confirmPromptValue := promptui.Prompt{
			Label: fmt.Sprintf("Repeat %s", promptLabel),
			Mask:  '*',
		}

		confirm, err := confirmPromptValue.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return "", errors.Errorf("failed read passphrase confirmation")
		}
		if secret != confirm {
			return "", errors.Errorf("%s do not match", promptLabel)
		}
	}

	return secret, nil
}
