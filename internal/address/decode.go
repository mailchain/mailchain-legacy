package address

import (
	"github.com/mailchain/mailchain/internal/chains"
	"github.com/pkg/errors"
)

func DecodeByProtocol(in, protocol string) ([]byte, error) {
	switch protocol {
	case chains.Ethereum:
		return decodeZeroX(in)
	default:
		return nil, errors.Errorf("")
	}
}
