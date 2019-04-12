package mailbox

import (
	"context"
)

// PubKeyFinder find public key to encrypt message with
type PubKeyFinder interface {
	PublicKeyFromAddress(ctx context.Context, network string, address []byte) ([]byte, error)
}
