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

package keystore

import (
	"github.com/mailchain/mailchain/internal/pkg/crypto/cipher"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
)

// Store private keys but does not return them, instead return decrypter or signer.
type Store interface {
	GetSigner(address []byte, chain string, deriveKeyOptions multi.OptionsBuilders) (mailbox.Signer, error)
	GetDecrypter(address []byte, decrypterType string, deriveKeyOptions multi.OptionsBuilders) (cipher.Decrypter, error)
	Store(private keys.PrivateKey, curveType string, deriveKeyOptions multi.OptionsBuilders) (address []byte, err error)
	HasAddress(address []byte) bool
	GetAddresses() ([][]byte, error)
}
