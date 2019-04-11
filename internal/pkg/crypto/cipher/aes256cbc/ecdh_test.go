package aes256cbc

import (
	"bytes"
	"encoding/hex"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
	"github.com/mailchain/mailchain/internal/pkg/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDerive(t *testing.T) {
	assert := assert.New(t)
	pub, err := secp256k1.PublicKeyToECIES(testutil.CharlottePublicKey)
	if err != nil {
		log.Fatal(err)
	}
	priv, err := secp256k1.PrivateKeyToECIES(testutil.SofiaPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	shared, err := deriveSharedSecret(pub, priv)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal("b6bdfade23178272425d25774a7d0d388fbef9480893fcc3646accc123eacc47", hex.EncodeToString(shared[:]))
}

func TestGenerateMacKeyAndEncryptionKey(t *testing.T) {
	assert := assert.New(t)
	secret, err := hex.DecodeString("04695325aac70f9f9ebe676248ebbfefa87b3eff16117559d2a0953d0e695be6")
	if err != nil {
		log.Fatal(err)
	}
	macKey, encryptionKey := generateMacKeyAndEncryptionKey(secret)

	assert.Equal("2cea25760305bdb3194057646bc46dc2eeee4890b711741c0b525454ac7c5ea8", hex.EncodeToString(macKey))
	assert.Equal("af0ad81e7d9194721d6c26f6c1f2a2b7fd06e2c99c4f5deefe59fb93936c981e", hex.EncodeToString(encryptionKey))
}

func TestGenerateIV(t *testing.T) {
	assert := assert.New(t)
	iv, err := generateIV()
	if err != nil {
		log.Fatal(err)
	}
	assert.Len(iv, 16)
}

func TestGenerateMac(t *testing.T) {
	assert := assert.New(t)
	macKey := testutil.MustHexDecodeString("2cea25760305bdb3194057646bc46dc2eeee4890b711741c0b525454ac7c5ea8")
	iv := testutil.MustHexDecodeString("05050505050505050505050505050505")
	cipherText := testutil.MustHexDecodeString("2ec66aac453ff543f47830d4b8cbc68d9965bf7c6bb69724fd4de26d41001256dfa6f7f0b3956ce21d4717caf75b0c2ad753852f216df6cfbcda4911619c5fc34798a19f81adff902c1ad906ab0edaec")
	tmpEphemeralPrivateKey, err := crypto.HexToECDSA("0404040404040404040404040404040404040404040404040404040404040404")
	if err != nil {
		log.Fatal(err)
	}
	ephemeralPrivateKey := ecies.ImportECDSA(tmpEphemeralPrivateKey)
	actual, err := generateMac(macKey, iv, ephemeralPrivateKey.PublicKey, cipherText)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal("4367ae8a54b65f99e4f2fd315ba65bf85e1138967a7bea451faf80f75cdf3404", hex.EncodeToString(actual))
}

func TestEncryptEncodeDecodeDecrypt(t *testing.T) {
	assert := assert.New(t)
	logrus.SetLevel(logrus.DebugLevel)
	cases := []struct {
		name                string
		recipientPublicKey  keys.PublicKey
		recipientPrivateKey keys.PrivateKey
		data                []byte
		err                 error
	}{
		{"to-sofia-short-text",
			testutil.SofiaPublicKey,
			testutil.SofiaPrivateKey,
			[]byte("Hi Sofia"),
			nil,
		},
		{"to-sofia-medium-text",
			testutil.SofiaPublicKey,
			testutil.SofiaPrivateKey,
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{"to-charlotte-short-text",
			testutil.CharlottePublicKey,
			testutil.CharlottePrivateKey,
			[]byte("Hi Charlotte"),
			nil,
		}, {"to-charlotte-medium-text",
			testutil.CharlottePublicKey,
			testutil.CharlottePrivateKey,
			[]byte("Hi Charlotte, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := Encrypt(tc.recipientPublicKey, tc.data)
			assert.Equal(tc.err, err)
			assert.NotNil(encrypted)
			encodedBytes, err := BytesEncode(*encrypted)
			assert.NoError(err)
			logrus.Debug(hex.EncodeToString(encodedBytes))
			toDecrypt, err := BytesDecode(encodedBytes)
			assert.NoError(err)

			toDecryptCopy := encryptedData{
				Ciphertext:                encrypted.Ciphertext,
				EphemeralPublicKey:        encrypted.EphemeralPublicKey,
				InitializationVector:      encrypted.InitializationVector,
				MessageAuthenticationCode: encrypted.MessageAuthenticationCode,
			}

			assert.Equal(encrypted, toDecrypt)
			assert.True(bytes.Equal(encrypted.Ciphertext, toDecryptCopy.Ciphertext))
			assert.True(bytes.Equal(encrypted.EphemeralPublicKey, toDecryptCopy.EphemeralPublicKey))
			assert.True(bytes.Equal(encrypted.InitializationVector, toDecryptCopy.InitializationVector))
			assert.True(bytes.Equal(encrypted.MessageAuthenticationCode, toDecryptCopy.MessageAuthenticationCode))

			assert.Equal(hex.EncodeToString(encrypted.Ciphertext), hex.EncodeToString(toDecryptCopy.Ciphertext))
			assert.Equal(hex.EncodeToString(encrypted.EphemeralPublicKey), hex.EncodeToString(toDecryptCopy.EphemeralPublicKey))
			assert.Equal(hex.EncodeToString(encrypted.InitializationVector), hex.EncodeToString(toDecryptCopy.InitializationVector))
			assert.Equal(hex.EncodeToString(encrypted.MessageAuthenticationCode), hex.EncodeToString(toDecryptCopy.MessageAuthenticationCode))

			nonEncodedDecrypt, err := decryptEncryptedData(tc.recipientPrivateKey, *encrypted)
			assert.Equal(tc.err, err)
			assert.Equal(tc.data, nonEncodedDecrypt)

			decryptedCopy, err := decryptEncryptedData(tc.recipientPrivateKey, toDecryptCopy)
			assert.Equal(tc.err, err)
			assert.Equal(tc.data, decryptedCopy)

			decrypted, err := decryptEncryptedData(tc.recipientPrivateKey, *toDecrypt)
			assert.Equal(tc.err, err)
			assert.Equal(tc.data, decrypted)

		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name                string
		recipientPublicKey  keys.PublicKey
		recipientPrivateKey keys.PrivateKey
		data                []byte
		err                 error
	}{
		{"to-sofia-short-text",
			testutil.SofiaPublicKey,
			testutil.SofiaPrivateKey,
			[]byte("Hi Sofia"),
			nil,
		}, {"to-sofia-medium-text",
			testutil.SofiaPublicKey,
			testutil.SofiaPrivateKey,
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
		{"to-charlotte-short-text",
			testutil.CharlottePublicKey,
			testutil.CharlottePrivateKey,
			[]byte("Hi Charlotte"),
			nil,
		}, {"to-charlotte-medium-text",
			testutil.CharlottePublicKey,
			testutil.CharlottePrivateKey,
			[]byte("Hi Charlotte, this is a little bit of a longer message to make sure there are no problems"),
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := Encrypt(tc.recipientPublicKey, tc.data)
			assert.Equal(tc.err, err)
			assert.NotNil(encrypted)

			decrypted, err := decryptEncryptedData(tc.recipientPrivateKey, *encrypted)
			assert.Equal(tc.err, err)
			assert.Equal(tc.data, decrypted)
		})
	}
}
