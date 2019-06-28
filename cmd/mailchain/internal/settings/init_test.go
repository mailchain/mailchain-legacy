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

package settings

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInitStore(t *testing.T) {
	type args struct {
		viper      *viper.Viper
		cfgFile    string
		logLevel   string
		createFile bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{viper.New(), "./testdata/.empty.yaml", "DEBUG", false},
			false,
		},
		{
			"err-invalid-file",
			args{viper.New(), "./testdata/.invalid.yaml", "DEBUG", false},
			true,
		},
		{
			"create-file",
			args{viper.New(),
				func() string {
					f := "./tmp/init/.create.yaml"
					os.RemoveAll("./tmp/init/")
					return f
				}(),
				"DEBUG", true},
			true,
		},
		{
			"invalid-level-empty-file",
			args{viper.New(), "", "INVALID", false},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitStore(tt.args.viper, tt.args.cfgFile, tt.args.logLevel, tt.args.createFile); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
