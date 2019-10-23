package aes256cbc

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name                string
		recipientPublicKey  crypto.PublicKey
		recipientPrivateKey crypto.PrivateKey
		data                []byte
		err                 error
	}{
		{
			"to-sofia-short-text",
			secp256k1test.SofiaPublicKey,
			secp256k1test.SofiaPrivateKey,
			[]byte("Hi Sofia"),
			nil,
		},
		{
			"to-sofia-medium-text",
			secp256k1test.SofiaPublicKey,
			secp256k1test.SofiaPrivateKey,
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{
			"to-charlotte-short-text",
			secp256k1test.CharlottePublicKey,
			secp256k1test.CharlottePrivateKey,
			[]byte("Hi Charlotte"),
			nil,
		},
		{
			"to-charlotte-medium-text",
			secp256k1test.CharlottePublicKey,
			secp256k1test.CharlottePrivateKey,
			[]byte("Hi Charlotte, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := NewEncrypter().Encrypt(tc.recipientPublicKey, tc.data)
			assert.Equal(tc.err, err)
			assert.NotNil(encrypted)
			decrypter := Decrypter{tc.recipientPrivateKey}

			decrypted, err := decrypter.Decrypt(encrypted)
			assert.Equal(tc.err, err)
			assert.Equal(tc.data, []byte(decrypted))
		})
	}
}
