package nacl

import (
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

const nonceSize = 24

func easyOpen(box, key []byte) ([]byte, error) {
	if len(box) < nonceSize {
		return nil, errors.New("secretbox: message too short")
	}
	decryptNonce := new([nonceSize]byte)
	copy(decryptNonce[:], box[:nonceSize])

	var secretKey [32]byte
	copy(secretKey[:], key)

	decrypted, ok := secretbox.Open([]byte{}, box[nonceSize:], decryptNonce, &secretKey)
	if !ok {
		return nil, errors.New("secretbox: could not decrypt data with private key")
	}
	return decrypted, nil
}

func easySeal(message, key []byte, rand io.Reader) ([]byte, error) {
	nonce := new([nonceSize]byte)
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	var secretKey [32]byte
	copy(secretKey[:], key)
	return secretbox.Seal(nonce[:], message, nonce, &secretKey), nil
}
