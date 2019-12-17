package datastore

import (
	"context"

	"github.com/mailchain/mailchain/crypto"
)

type PublicKeyStore interface {
	PutPublicKey(ctx context.Context, protocol, network string, address []byte, pubKey crypto.PublicKey) error
	GetPublicKey(ctx context.Context, protocol, network string, address []byte) (pubKey crypto.PublicKey, err error)
}
