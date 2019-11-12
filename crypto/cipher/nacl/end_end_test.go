package nacl

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/ed25519/ed25519test"
	"github.com/mailchain/mailchain/crypto/sr25519/sr25519test"
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
			ed25519test.SofiaPublicKey,
			ed25519test.SofiaPrivateKey,
			[]byte("Hi Sofia"),
			nil,
		},
		{
			"to-sofia-medium-text",
			ed25519test.SofiaPublicKey,
			ed25519test.SofiaPrivateKey,
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{
			"to-charlotte-short-text",
			ed25519test.CharlottePublicKey,
			ed25519test.CharlottePrivateKey,
			[]byte("Hi Charlotte"),
			nil,
		},
		{
			"to-charlotte-medium-text",
			ed25519test.CharlottePublicKey,
			ed25519test.CharlottePrivateKey,
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
