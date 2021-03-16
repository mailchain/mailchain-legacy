package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/mailchain/mailchain/sender"
	"github.com/pkg/errors"
)

// Send transaction using the RPC2 client.
func (c *Client) Send(ctx context.Context, network string, to, from, data []byte, txSigner signer.Signer, opts sender.SendOpts) error {
	algodClient, err := algod.MakeClient(c.algoAddress, c.algodToken)
	if err != nil {
		return errors.WithMessage(err, "failed to create algod client")
	}

	txParams, err := algodClient.BuildSuggestedParams()
	if err != nil {
		return errors.Wrap(err, "error getting suggested tx params")
	}

	fromAddress, _, err := addressing.EncodeByProtocol(from, protocols.Algorand)
	if err != nil {
		return err
	}

	toAddress, _, err := addressing.EncodeByProtocol(to, protocols.Algorand)
	if err != nil {
		return err
	}

	fromAddr := fromAddress
	toAddr := toAddress
	var amount uint64 = 0
	var minFee uint64 = 1000
	note := data
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash

	firstValidRound := uint64(txParams.FirstRoundValid)
	lastValidRound := uint64(txParams.LastRoundValid)

	logger := c.logger.With().Str("action", "send-message").Str("from", fromAddress).Str("to", toAddress).Logger()
	logger.Info().Msg("sending message")

	txn, err := transaction.MakePaymentTxn(fromAddr, toAddr, minFee, amount, firstValidRound, lastValidRound, note, "", genID, genHash)
	if err != nil {
		return errors.Wrap(err, "error creating transaction")
	}

	rawSignedTx, err := txSigner.Sign(algorand.SignerOptions{Transaction: txn})
	if err != nil {
		return errors.WithMessage(err, "could not sign transaction")
	}

	signedTx, ok := rawSignedTx.([]byte)
	if !ok {
		return errors.Errorf("sign did not return an []byte")
	}

	// Submit the transaction
	sendResponse, err := algodClient.SendRawTransaction(signedTx)
	if err != nil {
		return errors.WithMessage(err, "failed to send transaction")
	}

	logger.Info().Str("transaction-id", sendResponse.TxID).Msg("transaction submitted")
	return nil
}

func (c *Client) CheckBalance(address string, client *algod.Client) error {
	logger := c.logger.With().Str("action", "check-address").Logger()
	logger = logger.With().Str("address", address).Logger()

	accountInfo, err := client.AccountInformation(address)
	if err != nil {
		return err
	}
	// TODO: add fees to the check
	logger.Info().Msgf("account balance %d microAlgos", accountInfo.Amount)
	return nil
}
