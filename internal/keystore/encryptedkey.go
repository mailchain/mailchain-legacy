// Copyright 2022 Mailchain Ltd.
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
	"time"

	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
)

// EncryptedKey the data object when storing or retrieving a private key
type EncryptedKey struct {
	PublicKeyBytes []byte    `json:"public-key"`
	CipherText     []byte    `json:"cipher-text"`
	CurveType      string    `json:"curve-type"`
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	Version        string    `json:"version"`
	// KDF Key Definition Function
	StorageCipher string             `json:"storage-cipher"`
	KDF           string             `json:"kdf"`
	ScryptParams  *scrypt.DeriveOpts `json:"scrypt-params"`
}
