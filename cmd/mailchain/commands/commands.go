package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func exactAndOnlyValid(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			return err
		}
		return cobra.ExactArgs(n)(cmd, args)
	}
}

// formatExampleText pads the example text and adds a valid args section if supplied.
func formatExampleText(exampleText string, validArgs []string) string {
	if len(validArgs) == 0 {
		return fmt.Sprintf("  %s", exampleText)
	}
	return fmt.Sprintf("  %s\n\nValid arguments:\n  - %s", exampleText, strings.Join(validArgs, "\n  - "))
}
