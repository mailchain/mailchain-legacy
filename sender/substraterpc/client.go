package substraterpc

import (
	"context"
	"math/big"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mailchain/mailchain/internal/address"
)

//go:generate mockgen -source=client.go -package=substraterpctest -destination=./substraterpctest/client_mock.go

type Client interface {
	GetMetadata(blockHash types.Hash) (*types.Metadata, error)
	GetAddress(accountID []byte) types.Address
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	Call(metadata *types.Metadata, to types.Address, gas *big.Int, data []byte) (types.Call, error)
	NewExtrinsic(call types.Call) types.Extrinsic
	GetBlockHash(blockNumber uint64) (types.Hash, error)
	GetRuntimeVersion(blockHash types.Hash) (*types.RuntimeVersion, error)
	GetNonce(ctx context.Context, protocol, network string, addr []byte) (uint32, error)
	CreateSignatureOptions(blockHash, genesisHash types.Hash, mortalEra, immortalEra bool, rv types.RuntimeVersion, nonce, tip uint32) types.SignatureOptions
	SubmitExtrinsic(extrinsic *types.Extrinsic) (types.Hash, error)
}

const SuggestedGas = 32000

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

func (s SubstrateClient) GetAddress(accountID []byte) types.Address {
	return types.NewAddressFromAccountID(accountID)
}

func (s SubstrateClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(SuggestedGas), nil
}

func (s SubstrateClient) Call(metadata *types.Metadata, to types.Address, gas *big.Int, data []byte) (types.Call, error) {
	return types.NewCall(metadata, "Contracts.call", to, types.NewUCompactFromUInt(0), types.NewUCompactFromUInt(gas.Uint64()), data)
}

func (s SubstrateClient) NewExtrinsic(call types.Call) types.Extrinsic {
	return types.NewExtrinsic(call)
}

func (s SubstrateClient) GetBlockHash(blockNumber uint64) (types.Hash, error) {
	return s.api.RPC.Chain.GetBlockHash(blockNumber)
}

func (s SubstrateClient) GetRuntimeVersion(blockHash types.Hash) (*types.RuntimeVersion, error) {
	if (blockHash == types.Hash{}) {
		return s.api.RPC.State.GetRuntimeVersionLatest()
	}

	return s.api.RPC.State.GetRuntimeVersion(blockHash)
}

func (s SubstrateClient) GetNonce(ctx context.Context, protocol, network string, addr []byte) (uint32, error) {
	encodedAddress, _, err := address.EncodeByProtocol(addr, protocol)
	if err != nil {
		return uint32(0), err
	}

	idx, err := s.api.RPC.System.AccountNextIndex(encodedAddress)

	return uint32(idx), err
}

func (s SubstrateClient) CreateSignatureOptions(blockHash, genesisHash types.Hash, mortalEra, immortalEra bool, rv types.RuntimeVersion, nonce, tip uint32) types.SignatureOptions {
	return types.SignatureOptions{
		BlockHash:   blockHash,
		Era:         types.ExtrinsicEra{IsMortalEra: mortalEra, IsImmortalEra: immortalEra},
		GenesisHash: genesisHash,
		Nonce:       types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion: rv.SpecVersion,
		Tip:         types.NewUCompactFromUInt(uint64(tip)),
	}
}

func (s SubstrateClient) SubmitExtrinsic(extrinsic *types.Extrinsic) (types.Hash, error) {
	hash, err := s.api.RPC.Author.SubmitExtrinsic(*extrinsic)
	if err != nil {
		return types.Hash{}, err
	}

	return hash, nil
}
