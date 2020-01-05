package cipher

import "github.com/pkg/errors"

// DecodeBase58 returns the bytes represented by the base58 string src.
//
// DecodeBase58 expects that src contains only base58 characters.
// If the input is malformed, DecodeBase58 returns an error.
// func DecodeBase58(src string) ([]byte, error) {
// 	return base58.Decode(src)
// }

const (
	// ErrEncrypt returns the error message if encryption failed
	//
	ErrEncrypt = errors.New("cipher: encryption failed")
	// ErrDecrypt returns the error message if decryption failed
	//
	ErrDecrypt = errors.New("cipher: decryption failed")
)
