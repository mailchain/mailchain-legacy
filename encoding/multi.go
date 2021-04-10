package encoding

import (
	"errors"
	"strings"
)

var errUnsupportedEncoding = errors.New("encoding not supported")

// Decode returns the bytes represented by the decoded string src.
//
// Decode uses the decode method mapped to kind parameter.
// If the input is kind is unknown or the input is malformed for the decode method it returns an error.
func Decode(encoding, src string) ([]byte, error) {
	switch strings.ToLower(encoding) {
	case KindBase58:
		return DecodeBase58(src)
	case KindBase32:
		return DecodeBase32(src)
	case KindHex:
		return DecodeHex(src)
	case KindHex0XPrefix:
		return DecodeHexZeroX(src)
	case KindMnemonicAlgorand:
		return DecodeMnemonicAlgorand(src)
	default:
		return nil, errUnsupportedEncoding
	}
}
