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
	if err := s.Connect(); err != nil {
		return errors.WithMessage(err, "could not connect")
	}

	meta, err := s.client.GetMetadata(types.Hash{})
	if err != nil {
		return errors.WithMessage(err, "could not get latest metadata")
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.WithMessage(err, "could not get gas price")
	}

	addressTo := s.client.GetAddress(to[1:])

	c, err := s.client.Call(meta, addressTo, gasPrice, data)
	if err != nil {
		return errors.WithMessage(err, "could not create call")
	}

	ext := s.client.NewExtrinsic(c)

	genesisHash, err := s.client.GetBlockHash(0)
	if err != nil {
		return errors.WithMessage(err, "could not get block hash")
	}

	rv, err := s.client.GetRuntimeVersion(types.Hash{})
	if err != nil {
		return errors.WithMessage(err, "could not get runtime version")
	}

	nonce, err := s.client.GetNonce(ctx, protocols.Substrate, network, from)
	if err != nil {
		return errors.WithMessage(err, "could not get nonce")
	}

	o := s.client.CreateSignatureOptions(genesisHash, genesisHash, false, true, *rv, nonce, 0)

	signedExt, err := txSigner.Sign(substrate.SignerOptions{
		Extrinsic:        ext,
		SignatureOptions: o,
	})
	if err != nil {
		return errors.WithMessage(err, "could not sign the transaction")
	}

	signedExtTyped := signedExt.(*types.Extrinsic)

	hash, err := s.client.SubmitExtrinsic(signedExtTyped)
	if err != nil {
		return errors.WithMessage(err, "could not submit the transaction")
	}

	fmt.Printf("Transaction sent with hash %#x\n", hash)

	return err
}
