package nacl

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	cases := []struct {
		name                string
		recipientPublicKey  crypto.PublicKey
		recipientPrivateKey crypto.PrivateKey
		data                []byte
		err                 error
	}{
		{
			"to-alice-short-text",
			ed25519test.AlicePublicKey,
			ed25519test.AlicePrivateKey,
			[]byte("Hi Sofia"),
			nil,
		},
		{
			"to-alice-medium-text",
			ed25519test.AlicePublicKey,
			ed25519test.AlicePrivateKey,
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{
			"to-bob-short-text",
			ed25519test.BobPublicKey,
			ed25519test.BobPrivateKey,
			[]byte("Hi Charlotte"),
			nil,
		},
		{
			"to-bob-medium-text",
			ed25519test.BobPublicKey,
			ed25519test.BobPrivateKey,
			[]byte("Hi Charlotte, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{
			"to-alice-short-text",
			sr25519test.AlicePublicKey,
			sr25519test.AlicePrivateKey,
			[]byte("Hi Sofia"),
			nil,
		},
		{
			"to-alice-medium-text",
			sr25519test.AlicePublicKey,
			sr25519test.AlicePrivateKey,
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{
			"to-bob-short-text",
			sr25519test.BobPublicKey,
			sr25519test.BobPrivateKey,
			[]byte("Hi Charlotte"),
			nil,
		},
		{
			"to-bob-medium-text",
			sr25519test.BobPublicKey,
			sr25519test.BobPrivateKey,
			[]byte("Hi Charlotte, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{
			"to-alice-short-text",
			secp256k1test.AlicePublicKey,
			secp256k1test.AlicePrivateKey,
			[]byte("Hi Sofia"),
			nil,
		},
		{
			"to-alice-medium-text",
			secp256k1test.AlicePublicKey,
			secp256k1test.AlicePrivateKey,
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{
			"to-bob-short-text",
			secp256k1test.BobPublicKey,
			secp256k1test.BobPrivateKey,
			[]byte("Hi Charlotte"),
			nil,
		},
		{
			"to-bob-medium-text",
			secp256k1test.BobPublicKey,
			secp256k1test.BobPrivateKey,
			[]byte("Hi Charlotte, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encrypter, _ := NewEncrypter(tc.recipientPublicKey)
			encrypted, err := encrypter.Encrypt(tc.data)
			assert.Equal(t, tc.err, err)
			assert.NotNil(t, encrypted)

			decrypter, _ := NewDecrypter(tc.recipientPrivateKey)
			decrypted, err := decrypter.Decrypt(encrypted)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.data, []byte(decrypted))
		})
	}
}
