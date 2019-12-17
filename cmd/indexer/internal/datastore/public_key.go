package datastore

import (
	"github.com/mailchain/mailchain/crypto"
)

type PublicKeyStore interface {
	PutPublicKey(protocol, network string, address []byte, pubKey crypto.PublicKey) error
	GetPublicKey(protocol, network string, address []byte) (pubKey crypto.PublicKey, err error)
}
