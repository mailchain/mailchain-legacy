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

package aes256cbc

import "github.com/pkg/errors"

type encryptedData struct {
	InitializationVector      []byte `json:"iv"`
	EphemeralPublicKey        []byte `json:"ephemPublicKey"`
	Ciphertext                []byte `json:"ciphertext"`
	MessageAuthenticationCode []byte `json:"mac"`
}

func (e *encryptedData) verify() error {
	if len(e.InitializationVector) != 16 {
		return errors.Errorf("`InitializationVector` must be 16")
	}
	if len(e.EphemeralPublicKey) != 65 {
		return errors.Errorf("`EphemeralPublicKey` must be 65")
	}
	if len(e.MessageAuthenticationCode) != 32 {
		return errors.Errorf("`MessageAuthenticationCode` must be 16")
	}
	if len(e.Ciphertext) == 0 {
		return errors.Errorf("`Ciphertext` must not be empty")
	}

	return nil
}
