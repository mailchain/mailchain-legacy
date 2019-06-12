package commands

import (
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
