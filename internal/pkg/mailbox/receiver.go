package mailbox

import (
	"context"
	"errors"

	"github.com/mailchain/mailchain/internal/pkg/crypto/cipher"
)

// Receiver gets encrypted data from blockchain.
type Receiver interface {
	Receive(ctx context.Context, network string, address []byte) ([]cipher.EncryptedContent, error)
}

// type ReceiverOpts interface{}

var (
	errNetworkNotSupported = errors.New("network not supported")
)

// IsNetworkNotSupportedError network not supported errors can be resolved by selecting a different client or configuring the network.
func IsNetworkNotSupportedError(err error) bool {
	return err == errNetworkNotSupported
}
