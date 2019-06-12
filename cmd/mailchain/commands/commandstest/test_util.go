package commandstest

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func ExecuteCommandC(root *cobra.Command, args []string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	if err := root.ValidateArgs(args); err != nil {
		return nil, "", err
	}
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func AssertCommandOutput(t *testing.T, cmd *cobra.Command, err error, out string, wantOutput string) bool {
	assert := assert.New(t)
	if err == nil {
		if !assert.Equal(wantOutput, out) {
			t.Errorf("cmd().Execute().out = %v, want %v", out, wantOutput)
			return false
		}
	}
	if err != nil {
		if !assert.Equal(wantOutput+"\n"+cmd.UsageString()+"\n", out) {
			t.Errorf("cmd().Execute().out = %v, want %v", out, wantOutput)
			return false
		}
	}

	return true
}
