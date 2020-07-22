package substraterpc

import (
	"context"
	"fmt"

	"github.com/mailchain/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/mailchain/mailchain/sender"
	"github.com/pkg/errors"
)

func (s SubstrateRPC) Send(ctx context.Context, network string, to, from, data []byte, txSigner signer.Signer, opts sender.SendOpts) error {
	client := s.client

	meta, err := client.GetMetadata(types.Hash{})
	if err != nil {
		return errors.WithMessage(err, "could not get latest metadata")
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.WithMessage(err, "could not get gas price")
	}

	addressTo := client.GetAddress(to[1:])

	c, err := client.Call(meta, addressTo, gasPrice, data)
	if err != nil {
		return errors.WithMessage(err, "could not create call")
	}

	ext := client.NewExtrinsic(c)

	genesisHash, err := client.GetBlockHash(0)
	if err != nil {
		return errors.WithMessage(err, "could not get block hash")
	}

	rv, err := client.GetRuntimeVersion(types.Hash{})
	if err != nil {
		return errors.WithMessage(err, "could not get runtime version")
	}

	nonce, err := client.GetNonce(ctx, protocols.Substrate, network, from, meta)
	if err != nil {
		return errors.WithMessage(err, "could not get nonce")
	}

	o := client.CreateSignatureOptions(genesisHash, genesisHash, false, true, *rv, nonce, 0)

	signedExt, err := txSigner.Sign(substrate.SignerOptions{
		Extrinsic:        ext,
		SignatureOptions: o,
	})
	if err != nil {
		return errors.WithMessage(err, "could not sign the transaction")
	}

	signedExtTyped := signedExt.(*types.Extrinsic)

	hash, err := client.SubmitExtrinsic(signedExtTyped)
	if err != nil {
		return errors.WithMessage(err, "could not submit the transaction")
	}

	fmt.Printf("Transaction sent with hash %#x\n", hash)

	return err
}
