package substraterpc

import (
	"context"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"math/big"
)

//go:generate mockgen -source=client.go -package=substraterpctest -destination=./substraterpctest/client_mock.go

type Client interface {
	GetMetadata(blockHash types.Hash) (*types.Metadata, error)
	GetAddress(accountId []byte) types.Address
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	Call(metadata *types.Metadata, to types.Address, gas *big.Int, data []byte) (types.Call, error)
	NewExtrinsic(call types.Call) types.Extrinsic
	GetBlockHash(blockNumber uint64) (types.Hash, error)
	GetRuntimeVersion(blockHash types.Hash) (*types.RuntimeVersion, error)
	GetNonce(ctx context.Context, protocol, network string, address []byte, meta *types.Metadata) (uint32, error)
	CreateSignatureOptions(blockHash types.Hash, genesisHash types.Hash, mortalEra bool, nonce uint32, rv types.RuntimeVersion, tip uint32) types.SignatureOptions
	SubmitExtrinsic(extrinsic types.Extrinsic) (types.Hash, error)
}

type SubstrateClient struct {
	api *gsrpc.SubstrateAPI
}

func NewClient(api *gsrpc.SubstrateAPI) *SubstrateClient {
	return &SubstrateClient{api}
}

func (s SubstrateClient) GetMetadata(blockHash types.Hash) (*types.Metadata, error) {
	if (blockHash == types.Hash{}) {
		return s.api.RPC.State.GetMetadataLatest()
	}
	return s.api.RPC.State.GetMetadata(blockHash)
}

func (s SubstrateClient) GetAddress(accountId []byte) types.Address {
	return types.NewAddressFromAccountID(accountId)
}

func (s SubstrateClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(32000), nil
}

func (s SubstrateClient) Call(metadata *types.Metadata, to types.Address, gas *big.Int, data []byte) (types.Call, error) {
	return types.NewCall(metadata, "Contracts.call", to, types.UCompact(0), types.UCompact(gas.Uint64()), data)
}

func (s SubstrateClient) NewExtrinsic(call types.Call) types.Extrinsic {
	return types.NewExtrinsic(call)
}

func (s SubstrateClient) GetBlockHash(blockNumber uint64) (types.Hash, error) {
	if blockNumber == 0 {
		return s.api.RPC.Chain.GetBlockHashLatest()
	}
	return s.api.RPC.Chain.GetBlockHash(blockNumber)
}

func (s SubstrateClient) GetRuntimeVersion(blockHash types.Hash) (*types.RuntimeVersion, error) {
	if (blockHash == types.Hash{}) {
		return s.api.RPC.State.GetRuntimeVersionLatest()
	}
	return s.api.RPC.State.GetRuntimeVersion(blockHash)
}

func (s SubstrateClient) GetNonce(ctx context.Context, protocol, network string, address []byte, meta *types.Metadata) (uint32, error) {
	pkf := &substrate.PublicKeyFinder{}
	pk, err := pkf.PublicKeyFromAddress(ctx, protocol, network, address)
	if err != nil {
		return uint32(0), err
	}
	key, err := types.CreateStorageKey(meta, "System", "AccountNonce", pk.Bytes(), nil)
	if err != nil {
		return uint32(0), err
	}

	var nonce uint32
	err = s.api.RPC.State.GetStorageLatest(key, &nonce)
	return nonce, nil
}

func (s SubstrateClient) CreateSignatureOptions(blockHash types.Hash, genesisHash types.Hash, mortalEra bool, nonce uint32, rv types.RuntimeVersion, tip uint32) types.SignatureOptions {
	return types.SignatureOptions{
		BlockHash:   blockHash,
		Era:         types.ExtrinsicEra{IsMortalEra: mortalEra},
		GenesisHash: genesisHash,
		Nonce:       types.UCompact(nonce),
		SpecVersion: rv.SpecVersion,
		Tip:         types.UCompact(tip),
	}
}

func (s SubstrateClient) SubmitExtrinsic(extrinsic types.Extrinsic) (types.Hash, error) {
	hash, err := s.api.RPC.Author.SubmitExtrinsic(extrinsic)
	if err != nil {
		return types.Hash{}, err
	}
	return hash, nil
}
