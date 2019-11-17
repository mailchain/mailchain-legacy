package storage

import "github.com/mailchain/mailchain/internal/mail"

//go:generate mockgen -source=store.go -package=storagetest -destination=./storagetest/store_mock.go

// Store interface for using S3.
type Store interface {
	Exists(messageID mail.ID, contentsHash, integrityHash, contents []byte) error
	Put(messageID mail.ID, contentsHash, integrityHash, contents []byte) (address, resource string, mli uint64, err error)
}
