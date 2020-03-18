package s3

import (
	"context"
	"encoding/json"

	"github.com/mailchain/mailchain/encoding"

	"github.com/mailchain/mailchain/stores/s3store"
)

type TransactionStore struct {
	uploader s3store.Uploader
}

// NewS3TransactionStore creates a new S3 store.
func NewS3TransactionStore(region, bucket, id, secret string) (*TransactionStore, error) {
	s3Store, err := s3store.NewUploader(region, bucket, id, secret)
	if err != nil {
		return nil, err
	}

	return &TransactionStore{uploader: s3Store}, nil
}

// Key value of resource stored.
func (sts *TransactionStore) Key(hash []byte) string {
	return encoding.EncodeHex(hash)
}

type rawTransactionData struct {
	Protocol string      `json:"protocol"`
	Network  string      `json:"network"`
	Hash     string      `json:"hash"`
	Tx       interface{} `json:"transaction"`
}

func (sts *TransactionStore) PutRawTransaction(ctx context.Context, protocol, network string, hash []byte, tx interface{}) error {
	rtd := rawTransactionData{
		Protocol: protocol,
		Network:  network,
		Hash:     encoding.EncodeHexZeroX(hash),
		Tx:       tx,
	}

	jsonBody, err := json.Marshal(rtd)
	if err != nil {
		return err
	}

	key := sts.Key(hash)
	_, err = sts.uploader.Upload(ctx, nil, key, jsonBody)

	return err
}
