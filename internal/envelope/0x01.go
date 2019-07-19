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
			string(hash),
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
