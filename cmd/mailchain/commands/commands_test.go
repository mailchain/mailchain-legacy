package commands

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_exactAndOnlyValid(t *testing.T) {
	type args struct {
		n int
	}
	type cmd struct {
		args      []string
		validArgs []string
	}
	tests := []struct {
		name    string
		args    args
		cmd     cmd
		wantErr bool
	}{
		{
			"success",
			args{
				1,
			},
			cmd{
				[]string{"mailchain"},
				[]string{"s3", "mailchain"},
			},
			false,
		},
		{
			"err-invalid-arg",
			args{
				1,
			},
			cmd{
				[]string{"invalid"},
				[]string{"s3", "mailchain"},
			},
			true,
		},
		{
			"err-too-many-args",
			args{
				1,
			},
			cmd{
				[]string{"mailchain", "s3"},
				[]string{"s3", "mailchain"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exactAndOnlyValid(tt.args.n)
			c := &cobra.Command{
				ValidArgs: tt.cmd.validArgs,
			}
			err := got(c, tt.cmd.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("exactAndOnlyValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_formatExampleText(t *testing.T) {

	type args struct {
		exampleText string
		validArgs   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"no-args",
			args{
				"this is example text",
				nil,
			},
			"  this is example text",
		},
		{
			"single-arg",
			args{
				"this is example text",
				[]string{"arg1"},
			},
			`  this is example text

Valid arguments:
  - arg1`,
		},
		{
			"multi-arg",
			args{
				"this is example text",
				[]string{"arg1", "arg2"},
			},
			`  this is example text

Valid arguments:
  - arg1
  - arg2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExampleText(tt.args.exampleText, tt.args.validArgs)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("formatExampleText() = %v, want %v", got, tt.want)
			}
		})
	}
}
