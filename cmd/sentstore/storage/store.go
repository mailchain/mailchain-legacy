package storage

import "github.com/mailchain/mailchain/internal/mail"

//go:generate mockgen -source=store.go -package=storagetest -destination=./storagetest/store_mock.go
type Store interface {
	Exists(messageID mail.ID, contents []byte, hash string) error
	Put(messageID mail.ID, contents []byte, hash string) (string, error)
}
