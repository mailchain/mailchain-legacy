package s3

import (
	"context"
	"encoding/json"

	"github.com/mailchain/mailchain/stores/s3store"
)

type S3TransactionStore struct {
	*s3store.S3Store
}

// NewS3TransactionStore creates a new S3 store.
func NewS3TransactionStore(region, bucket, id, secret string) (*S3TransactionStore, error) {
	s3Store, err := s3store.NewS3Store(region, bucket, id, secret)
	if err != nil {
		return nil, err
	}
	return &S3TransactionStore{S3Store: s3Store}, nil
}

// Key value of resource stored.
func (sts *S3TransactionStore) Key(hash []byte) string {
	return sts.EncodeKey(hash)
}

type rawTransactionData struct {
	Protocol string      `json:"protocol"`
	Network  string      `json:"network"`
	Hash     []byte      `json:"hash"`
	Tx       interface{} `json:"transaction"`
}

func (sts *S3TransactionStore) PutRawTransaction(ctx context.Context, protocol, network string, hash []byte, tx interface{}) error {
	rtd := rawTransactionData{
		Protocol: protocol,
		Network:  network,
		Hash:     hash,
		Tx:       tx,
	}

	jsonBody, err := json.Marshal(rtd)
	if err != nil {
		return err
	}

	key := sts.EncodeKey(hash)
	_, _, _, err = sts.Upload(ctx, nil, key, jsonBody)
	return err
}
