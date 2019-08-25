package address

import (
	"github.com/mailchain/mailchain/internal/chains"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/pkg/errors"
)

func EncodeByProtocol(in []byte, protocol string) (string, error) {
	switch protocol {
	case chains.Ethereum:
		return encoding.EncodeZeroX(in), nil
	default:
		return "", errors.Errorf("%q unsupported protocol", protocol)
	}
}
