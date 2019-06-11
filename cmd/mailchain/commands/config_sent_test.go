// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_configStorageSent(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		viper        *viper.Viper
		sentSelector setup.SimpleSelector
	}
	tests := []struct {
		name        string
		args        args
		cmdArgs     []string
		wantOutput  string
		wantExecErr bool
	}{
		{
			"success",
			args{
				viper.New(),
				nil,
			},
			[]string{},
			"out",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfgStorageSent(tt.args.viper, tt.args.sentSelector)
			if !assert.NotNil(got) {
				t.Error("cfgStorageSent() is nil")
			}
			// cobra.ap
			_, out, err := executeCommandC(got, tt.cmdArgs)
			if (err != nil) != tt.wantExecErr {
				t.Errorf("cfgStorageSent().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
				return
			}
			if !assert.Equal(tt.wantOutput, out) {
				t.Errorf("cfgStorageSent().Execute().out = %v, want %v", out, tt.wantOutput)
			}
			// viper.GetViper()
			// viper.
			viper.Reset()
			// if ; !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("cfgStorageSent() = %v, want %v", got, tt.want)
			// }
		})
	}
}
func executeCommandC(root *cobra.Command, args []string) (c *cobra.Command, output string, err error) {
	// root.
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
