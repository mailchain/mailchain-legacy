package commandstest

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// ExecuteCommandC will run a command to capture the output.
func ExecuteCommandC(root *cobra.Command, args []string, flags map[string]string) (c *cobra.Command, output string, err error) {
	if err := root.ValidateArgs(args); err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)

	root.SetOutput(buf)
	root.SetArgs(args)

	for x := range flags {
		_ = root.Flags().Set(x, flags[x])
	}

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

// AssertCommandOutput ensure that the command outputs a specific string.
func AssertCommandOutput(t *testing.T, cmd *cobra.Command, err error, out, wantOutput string) bool {
	if err == nil {
		if !assert.Equal(t, wantOutput, out) {
			t.Errorf("cmd().Execute().out = %v, want %v", out, wantOutput)
			return false
		}
	}

	if err != nil {
		if !assert.Equal(t, wantOutput+"\n"+cmd.UsageString()+"\n", out) {
			t.Errorf("cmd().Execute().out = %v, want %v", out, wantOutput)
			return false
		}
	}

	return true
}
