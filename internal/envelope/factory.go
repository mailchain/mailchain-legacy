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

package envelope

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

// NewEnvelope returns a new envelope
func NewEnvelope(encrypter cipher.Encrypter, pubkey crypto.PublicKey, o []CreateOptionsBuilder) (Data, error) {
	opts := &CreateOpts{}
	apply(opts, o)

	switch opts.Kind {
	case Kind0x01:
		return NewZeroX01(encrypter, pubkey, opts)
	default:
		return nil, errors.Errorf("unknown kind")
	}
}

func apply(o *CreateOpts, opts []CreateOptionsBuilder) {
	for _, f := range opts {
		f(o)
	}
}

// ParseEnvelope parses envelope from string to byte
func ParseEnvelope(envelope string) (byte, error) {
	switch envelope {
	case KindString0x01:
		return Kind0x01, nil
	case KindString0x50:
		return Kind0x50, nil
	default:
		return 0x0, errors.Errorf("`envelope` provided is invalid")
	}
}
