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
	"net/url"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/internal/hash"
	"github.com/mailchain/mailchain/internal/mli"
	"github.com/pkg/errors"
)

// NewZeroX02 creates a new envelope of type ZeroX02.
// ZeroX02 envelope allows sending private messages with the minimal bytes by using `Uint64Bytes` where encryptedHash is the location.
func NewZeroX02(encrypter cipher.Encrypter, opts *CreateOpts) (*ZeroX02, error) {
	if opts.Location == 0 {
		return nil, errors.Errorf("location must be set")
	}

	if len(opts.EncryptedContents) == 0 {
		return nil, errors.Errorf("EncryptedContents must not be empty")
	}

	if opts.Resource == "" {
		return nil, errors.Errorf("resource must not be empty")
	}

	dec, err := cid.Decode(opts.Resource)
	if err != nil {
		return nil, errors.WithMessage(err, "0x02: cid decode failed")
	}

	cidHash, err := hash.Create(hash.CIVv1SHA2256Raw, opts.EncryptedContents)
	if err != nil {
		return nil, errors.WithMessage(err, "0x02: cid hash could not be created")
	}

	if !bytes.Equal(cidHash, dec.Bytes()) {
		return nil, errors.Errorf("encrypted contents and resource hash must match")
	}

	locHash := NewUInt64Bytes(opts.Location, dec.Bytes())

	enc, err := encrypter.Encrypt(cipher.PlainContent(locHash))
	if err != nil {
		return nil, err
	}

	env := &ZeroX02{
		UIBEncryptedLocationHash: enc,
		DecryptedHash:            opts.DecryptedHash,
	}

	return env, nil
}

// URL returns the addressable location of the message, the URL may be encrypted requiring decrypter to be supplied.
// URL is contained in the UIBEncryptedLocationHash which must first be decrypted.
// The decrypted data is converted to `UInt64Bytes`. The extracted identified is used to look up the Message Location Indicator (MLI).
// MLI address and hash are combined to make an addressable URL.
func (x *ZeroX02) URL(decrypter cipher.Decrypter) (*url.URL, error) {
	decrypted, err := decrypter.Decrypt(x.UIBEncryptedLocationHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	locationHash := UInt64Bytes(decrypted)

	code, locHash, err := locationHash.Values()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	loc, ok := mli.ToAddress()[code]
	if !ok {
		return nil, errors.Errorf("unknown location code %q", code)
	}

	decoded, err := cid.Cast(locHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return url.Parse(strings.Join(
		[]string{loc, decoded.String()},
		"/"))
}

// ContentsHash returns a hash of the decrypted content.
// This can be used to verify the contents of the message have not been tampered with.
// Returns the value stored in DecryptedHash.
func (x *ZeroX02) ContentsHash(decrypter cipher.Decrypter) ([]byte, error) {
	return x.DecryptedHash, nil
}

// IntegrityHash returns a hash of the encrypted content. This can be used to validate the integrity of the contents before decrypting.
// UIBEncryptedLocationHash is decrypted to get a location hash. This is a UInt64Bytes and the data portion is the value for IntegrityHash.
func (x *ZeroX02) IntegrityHash(decrypter cipher.Decrypter) ([]byte, error) {
	decrypted, err := decrypter.Decrypt(x.UIBEncryptedLocationHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	locationHash := UInt64Bytes(decrypted)

	return locationHash.Bytes()
}

// Valid checks the envelopes contents for no integrity issues which would prevent the envelope from being read.
func (x *ZeroX02) Valid() error {
	if len(x.UIBEncryptedLocationHash) == 0 {
		return errors.Errorf("`EncryptedLocationHash` must not be empty")
	}

	if len(x.DecryptedHash) == 0 {
		return errors.Errorf("`DecryptedHash` must not be empty")
	}

	return nil
}

func (x *ZeroX02) DecrypterKind() (byte, error) {
	if len(x.UIBEncryptedLocationHash) == 0 {
		return 0x0, errors.Errorf("`EncryptedLocationHash` must not be empty")
	}

	return x.UIBEncryptedLocationHash[0], nil
}
