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

import "golang.org/x/crypto/scrypt"

// DeriveKey from the options provided.
func DeriveKey(o []DeriveOptionsBuilder) ([]byte, error) {
	opts := &DeriveOpts{}
	apply(opts, o)
	return scrypt.Key([]byte(opts.Passphrase), opts.Salt, opts.N, opts.R, opts.P, opts.Len)
}

// CreateOptions that can be stored.
func CreateOptions(o []DeriveOptionsBuilder) *DeriveOpts {
	opts := &DeriveOpts{}
	apply(opts, o)
	return opts
}

func apply(o *DeriveOpts, opts []DeriveOptionsBuilder) {
	for _, f := range opts {
		f(o)
	}
}
