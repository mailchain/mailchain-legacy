package encoding

import (
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
)

func EncodeZeroX(in []byte) string {
	out := make([]byte, len(in)*2+2)
	copy(out, "0x")
	hex.Encode(out[2:], in)
	return string(out)
}

func DecodeZeroX(in string) ([]byte, error) {
	if in == "" {
		return nil, errors.Errorf("empty hex string")
	}
	if !strings.HasPrefix(in, "0x") {
		return nil, errors.Errorf("missing \"0x\" prefix from hex string")
	}
	return hex.DecodeString(in[2:])
}
