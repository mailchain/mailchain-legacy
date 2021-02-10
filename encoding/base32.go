package encoding

import "encoding/base32"

// DecodeBase32 returns the bytes represented by the base32 string src.
//
// DecodeBase32 expects that src contains only base32 characters.
// If the input is malformed, DecodeBase32 returns an error.
func DecodeBase32(src string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(src)
}

// EncodeBase32 returns the string represented by the base32 byte src.
//
// EncodeBase32 expects that src contains only base32 byte.
// If the input is malformed, EncodeBase32 returns an error.
func EncodeBase32(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}
