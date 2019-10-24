package nacl

import (
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

const nonceSize = 24
const secretKeySize = 32

func easyOpen(box, key []byte) ([]byte, error) {
	if len(key) != secretKeySize {
		return nil, errors.New("secretbox: key length must be 32")
	}
	if len(box) < nonceSize {
		return nil, errors.New("secretbox: message too short")
	}
	decryptNonce := new([nonceSize]byte)
	copy(decryptNonce[:], box[:nonceSize])

	var secretKey [secretKeySize]byte

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

	var secretKey [secretKeySize]byte

	copy(secretKey[:], key)

	return secretbox.Seal(nonce[:], message, nonce, &secretKey), nil
}
