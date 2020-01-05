package cipher

import (
	"github.com/pkg/errors"
)

// ErrEncrypt returns the error message if encryption failed
//
func ErrEncrypt() error {
	return errors.New("cipher: encryption failed")
}

// ErrDecrypt returns the error message if decryption failed
//
func ErrDecrypt() error {
	return errors.New("cipher: decryption failed")
}
