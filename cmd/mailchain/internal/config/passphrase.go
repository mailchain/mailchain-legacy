package config

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra" // nolint: depguard
	"github.com/ttacon/chalk"
)

// Passphrase is extracted from the command
func Passphrase(cmd *cobra.Command) (string, error) {
	passphrase, err := cmd.Flags().GetString("passphrase")
	if err != nil {
		return "", errors.WithMessage(err, "could not get `passphrase`")
	}
	if passphrase != "" {
		return passphrase, nil
	}
	emptyPassphrase, err := cmd.Flags().GetBool("empty-passphrase")
	if err != nil {
		return "", errors.WithMessage(err, "could not get `empty-passphrase`")
	}
	if emptyPassphrase {
		return "", nil
	}
	fmt.Println(chalk.Yellow, "Note: To derive a storage key passphrase is required. The passphrase must be secure and not guessable.")
	return passphraseFromPrompt()
}

func passphraseFromPrompt() (string, error) {
	prompt := promptui.Prompt{
		Label: "Passphrase",
		Mask:  '*',
	}
	password, err := prompt.Run()
	if err != nil {
		return "", errors.Errorf("failed read passphrase")
	}

	confirmPrompt := promptui.Prompt{
		Label: "Repeat passphrase: ",
		Mask:  '*',
	}
	confirm, err := confirmPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", errors.Errorf("failed read passphrase confirmation")
	}
	if password != confirm {
		return "", errors.Errorf("Passphrases do not match")
	}

	return password, nil
}
