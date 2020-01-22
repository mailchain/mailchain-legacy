package os

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/encoding"
	"github.com/spf13/afero"
)

// NewRawTransactionStore creates a new raw transaction store with the path specified
func NewRawTransactionStore(path string) (datastore.RawTransactionStore, error) {
	return &RawTransactionStore{fs: afero.NewBasePathFs(afero.NewOsFs(), path)}, nil
}

// RawTransactionStore object
type RawTransactionStore struct {
	fs afero.Fs
}

type rawTransaction struct {
	Protocol string      `json:"protocol"`
	Network  string      `json:"network"`
	Hash     []byte      `json:"hash"`
	Tx       interface{} `json:"transaction"`
}

// PutRawTransaction writes the raw transaction to the file system
func (s RawTransactionStore) PutRawTransaction(ctx context.Context, protocol, network string, hash []byte, tx interface{}) error {
	rawTransactionJSON := rawTransaction{
		Protocol: protocol,
		Network:  network,
		Hash:     hash,
		Tx:       tx,
	}

	// This cannot fail here, as the only possible failures would be:
	// passing an invalid type like a channel or
	// passing an invalid value like math.Inf
	rawTransaction, _ := json.Marshal(rawTransactionJSON)

	fileName := fmt.Sprintf("%s.json", encoding.EncodeHex(rawTransaction))

	const filePerm = 0700
	err := afero.WriteFile(s.fs, fileName, rawTransaction, filePerm)

	if err != nil {
		return err
	}

	return nil
}
