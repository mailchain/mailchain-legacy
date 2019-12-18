package datastore

import (
	"context"

	"github.com/mailchain/mailchain/crypto"
)

//go:generate mockgen -source=public_key.go -package=datastoretest -destination=./datastoretest/publick_key_mock.go


type PublicKeyStore interface {
	PutPublicKey(ctx context.Context, protocol, network string, address []byte, pubKey crypto.PublicKey) error
	GetPublicKey(ctx context.Context, protocol, network string, address []byte) (pubKey crypto.PublicKey, err error)
}
