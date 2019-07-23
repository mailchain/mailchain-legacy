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
	"bytes"
	"encoding/hex"
	"net/url"
	"strings"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/pkg/errors"
)

func NewZeroX01(encrypter cipher.Encrypter, pubkey crypto.PublicKey, opts *CreateOpts) (*ZeroX01, error) {
	if opts.Location == 0 {
		return nil, errors.Errorf("location must be set")
	}
	if len(opts.DecryptedHash) == 0 {
		return nil, errors.Errorf("decryptedHash must not be empty")
	}
	if opts.Resource == "" {
		return nil, errors.Errorf("resource must not be empty")
	}
	resource, err := hex.DecodeString(opts.Resource)
	if err != nil {
		return nil, errors.Errorf("resource could not be decoded")
	}
	if !bytes.Equal(resource, opts.DecryptedHash) {
		return nil, errors.Errorf("resource %q and decrypted hash %q must match",
			hex.EncodeToString(resource), hex.EncodeToString(opts.DecryptedHash))
	}

	locHash := NewUInt64Bytes(opts.Location, opts.DecryptedHash)

	enc, err := encrypter.Encrypt(pubkey, cipher.PlainContent(locHash))
	if err != nil {
		return nil, err
	}

	env := &ZeroX01{
		UIBEncryptedLocationHash: enc,
		EncryptedHash:            opts.EncryptedHash,
	}
	return env, nil
}

func (d *ZeroX01) URL(decrypter cipher.Decrypter) (*url.URL, error) {
	decrypted, err := decrypter.Decrypt(d.UIBEncryptedLocationHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	locationHash := UInt64Bytes(decrypted)

	code, hash, err := locationHash.Values()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	loc, ok := MLIToAddress()[code]
	if !ok {
		return nil, errors.Errorf("unknown location code %q", code)
	}
	return url.Parse(strings.Join(
		[]string{
			loc,
			hex.EncodeToString(hash),
		},
		"/"))
}

func (d *ZeroX01) ContentsHash(decrypter cipher.Decrypter) ([]byte, error) {
	decrypted, err := decrypter.Decrypt(d.UIBEncryptedLocationHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	locationHash := UInt64Bytes(decrypted)

	return locationHash.Bytes()
	// TODO: validate hash
}

func (d *ZeroX01) IntegrityHash(decrypter cipher.Decrypter) ([]byte, error) {
	return d.EncryptedHash, nil
}

func (d *ZeroX01) Valid() error {
	if len(d.UIBEncryptedLocationHash) == 0 {
		return errors.Errorf("`EncryptedLocationHash` must not be empty")
	}

	return nil
}
