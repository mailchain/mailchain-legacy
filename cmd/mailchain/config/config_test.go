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

package config

import "testing"

func TestInit(t *testing.T) {
	type args struct {
		cfgFile  string
		logLevel string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{"./testdata/.empty.yaml", "DEBUG"},
			false,
		},
		{
			"err-invalid-file",
			args{"./testdata/.invalid.yaml", "DEBUG"},
			true,
		},
		{
			"err-no-file",
			args{"./testdata/.no-file.yaml", "DEBUG"},
			true,
		},
		{
			"invalid-level",
			args{"./testdata/.empty.yaml", "INVALID"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Init(tt.args.cfgFile, tt.args.logLevel); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
