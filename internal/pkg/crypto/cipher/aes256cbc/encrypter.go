package aes256cbc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/andreburgaud/crypt2go/padding"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys"
	"github.com/mailchain/mailchain/internal/pkg/crypto/keys/secp256k1"
	"github.com/pkg/errors"
)

// Encrypt data using recipient public key with AES in CBC mode.  Generate an ephemeral private key and IV.
func Encrypt(recipientPublicKey keys.PublicKey, message []byte) (*encryptedData, error) {
	rpk, err := secp256k1.PublicKeyToECIES(recipientPublicKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert pk")
	}

	ephemeral, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate ephemeral key")
	}
	iv, err := generateIV()
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate iv")
	}
	return encrypt(ephemeral, rpk, message, iv)
}

func encrypt(ephemeralPrivateKey *ecies.PrivateKey, pub *ecies.PublicKey, input []byte, iv []byte) (*encryptedData, error) {
	ephemeralPublicKey := ephemeralPrivateKey.PublicKey
	sharedSecret, err := deriveSharedSecret(pub, ephemeralPrivateKey)
	if err != nil {
		return nil, err
	}
	macKey, encryptionKey := generateMacKeyAndEncryptionKey(sharedSecret)
	ciphertext, err := encryptCBC(input, iv, encryptionKey)
	if err != nil {
		return nil, errors.WithMessage(err, "encryptCBC failed")
	}

	mac, err := generateMac(macKey, iv, ephemeralPublicKey, ciphertext)
	if err != nil {
		return nil, errors.WithMessage(err, "generateMac failed")
	}

	return &encryptedData{
		MessageAuthenticationCode: mac,
		InitializationVector:      iv,
		EphemeralPublicKey:        elliptic.Marshal(ecies.DefaultCurve, ephemeralPublicKey.X, ephemeralPublicKey.Y),
		Ciphertext:                ciphertext,
	}, nil
}

func encryptCBC(data, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data, err = padding.NewPkcs7Padding(block.BlockSize()).Pad(data)
	if err != nil {
		return nil, errors.WithMessage(err, "could not pad")
	}

	ciphertext := make([]byte, len(data))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, data)

	return ciphertext, nil
}
