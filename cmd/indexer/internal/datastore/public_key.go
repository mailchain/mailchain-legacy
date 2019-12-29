package datastore

import (
	"context"

	"github.com/mailchain/mailchain/crypto"
)

//go:generate mockgen -source=public_key.go -package=datastoretest -destination=./datastoretest/publick_key_mock.go

type PublicKey struct {
	PublicKey crypto.PublicKey
	BlockHash []byte
	TxHash    []byte
}

type PublicKeyStore interface {
	PutPublicKey(ctx context.Context, protocol, network string, address []byte, pubKey *PublicKey) error
	GetPublicKey(ctx context.Context, protocol, network string, address []byte) (pubKey *PublicKey, err error)
}
