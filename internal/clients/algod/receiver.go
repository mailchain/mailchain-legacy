package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/stores"
	"github.com/pkg/errors"
)

// Receive check algorand transactions for mailchain messages.
func (c *Client) Receive(ctx context.Context, protocol, network string, address []byte) ([]stores.Transaction, error) {
	if !c.isNetworkSupported(network) {
		return nil, errors.New("network not supported")
	}

	algodClient, err := indexer.MakeClient(c.networkConfigs[network].url, c.algodToken)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create algod client")
	}

	txResult, err := algodClient.SearchForTransactions().NotePrefix(encoding.DataPrefix()).AddressRole("receiver").Address(toAlgodAddress(address)).Do(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := []stores.Transaction{}

	for i := range txResult.Transactions {
		tx := txResult.Transactions[i]

		res = append(res, stores.Transaction{
			EnvelopeData: tx.Note[len(encoding.DataPrefix()):],
			BlockNumber:  int64(tx.RoundTime),
			Hash:         []byte(tx.Id),
		})
	}

	return res, nil
}

func (c *Client) Kind() string {
	return "algod"
}

func toAlgodAddress(address []byte) types.Address {
	var out types.Address

	copy(out[:], address)

	return out
}
