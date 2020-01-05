package cipher

import "errors"

var (
	// ErrEncrypt returns the error message if encryption failed
	//
	ErrEncrypt = errors.New("cipher: encryption failed") //nolint:gochecknoglobals
	// ErrDecrypt returns the error message if decryption failed
	//
	ErrDecrypt = errors.New("cipher: decryption failed") //nolint:gochecknoglobals
)
