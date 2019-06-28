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

// import (
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/mailchain/mailchain/cmd/mailchain/commands/commandstest"
// 	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup"
// 	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup/setuptest"
// 	"github.com/pkg/errors"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_configStorageSent(t *testing.T) {
// 	assert := assert.New(t)
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
// 	type args struct {
// 		sentSelector setup.SimpleSelector
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		cmdArgs     []string
// 		cmdFlags    map[string]string
// 		wantOutput  string
// 		wantExecErr bool
// 	}{
// 		{
// 			"success",
// 			args{
// 				func() setup.SimpleSelector {
// 					g := setuptest.NewMockSimpleSelector(mockCtrl)
// 					g.EXPECT().Select("mailchain").Return("mailchain", nil)
// 					return g
// 				}(),
// 			},
// 			[]string{"mailchain"},
// 			nil,
// 			"Sent store \"mailchain\" configured\n",
// 			false,
// 		},
// 		{
// 			"err-selector",
// 			args{
// 				func() setup.SimpleSelector {
// 					g := setuptest.NewMockSimpleSelector(mockCtrl)
// 					g.EXPECT().Select("mailchain").Return("", errors.Errorf("selector failed"))
// 					return g
// 				}(),
// 			},
// 			[]string{"mailchain"},
// 			nil,
// 			"Error: selector failed",
// 			true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := configStorageSent(tt.args.sentSelector)
// 			if !assert.NotNil(got) {
// 				t.Error("cfgStorageSent() is nil")
// 			}
// 			_, out, err := commandstest.ExecuteCommandC(got, tt.cmdArgs, tt.cmdFlags)
// 			if (err != nil) != tt.wantExecErr {
// 				t.Errorf("cfgStorageSent().execute() error = %v, wantExecErr %v", err, tt.wantExecErr)
// 				return
// 			}
// 			if !commandstest.AssertCommandOutput(t, got, err, out, tt.wantOutput) {
// 				t.Errorf("cfgStorageSent().Execute().out != %v", tt.wantOutput)
// 			}
// 		})
// 	}
// }
