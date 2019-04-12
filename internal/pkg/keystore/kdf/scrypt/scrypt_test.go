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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDefault(t *testing.T) {
	assert := assert.New(t)
	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions()})
	assert.Equal(32, opts.Len)
	assert.Equal(262144, opts.N)
	assert.Equal(1, opts.P)
	assert.Equal("", opts.Passphrase)
	assert.Equal(8, opts.R)
	assert.Nil(opts.Salt)
}

func TestApplyDefaultAndPassword(t *testing.T) {
	assert := assert.New(t)
	randomSalt, err := RandomSalt()
	if err != nil {
		t.Fail()
	}

	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions(), WithPassphrase("test"), randomSalt})
	assert.Equal(32, opts.Len)
	assert.Equal(262144, opts.N)
	assert.Equal(1, opts.P)
	assert.Equal("test", opts.Passphrase)
	assert.Equal(8, opts.R)
	assert.Equal(32, len(opts.Salt))
}
