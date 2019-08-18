package address

import (
	"github.com/mailchain/mailchain/internal/chains"
	"github.com/pkg/errors"
)

func EncodeByProtocol(in []byte, protocol string) (string, error) {
	switch protocol {
	case chains.Ethereum:
		return encodeZeroX(in), nil
	default:
		return "", errors.Errorf("%q unsupported protocol", protocol)
	}
}
