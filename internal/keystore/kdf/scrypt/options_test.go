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

package scrypt

import "testing"

func TestDeriveOpts_KDF(t *testing.T) {
	type fields struct {
		Len        int
		N          int
		P          int
		R          int
		Salt       []byte
		Passphrase string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"success",
			fields{},
			"scrypt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DeriveOpts{
				Len:        tt.fields.Len,
				N:          tt.fields.N,
				P:          tt.fields.P,
				R:          tt.fields.R,
				Salt:       tt.fields.Salt,
				Passphrase: tt.fields.Passphrase,
			}
			if got := d.KDF(); got != tt.want {
				t.Errorf("DeriveOpts.KDF() = %v, want %v", got, tt.want)
			}
		})
	}
}
