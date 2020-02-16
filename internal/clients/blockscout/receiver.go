package blockscout

import (
	"bytes"
	"context"

	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

// Receive check ethereum transactions for mailchain messages
func (c APIClient) Receive(ctx context.Context, network string, address []byte) ([]mailbox.Transaction, error) {
	if !c.isNetworkSupported(network) {
		return nil, errors.Errorf("network not supported")
	}
	txResult, err := c.getTransactionsByAddress(network, address)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := []mailbox.Transaction{}
	txHashes := map[string]bool{}

	for i := range txResult.Result {
		x := txResult.Result[i]

		_, ok := txHashes[x.Hash]
		if ok {
			continue
		}
		txHashes[x.Hash] = true
		encryptedTransactionData, err := encoding.DecodeHexZeroX(x.Input)
		if err != nil {
			continue // invalid data should move to next record
		}

		if !bytes.HasPrefix(encryptedTransactionData, encoding.DataPrefix()) {
			continue
		}

		res = append(res, mailbox.Transaction{
			Data:    encryptedTransactionData[len(encoding.DataPrefix()):],
			BlockID: []byte(x.BlockNumber),
			Hash:    []byte(x.Hash),
		})
	}
	return res, nil
}
