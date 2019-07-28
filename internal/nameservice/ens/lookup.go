package ens

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mailchain/mailchain/internal/nameservice"
	ens "github.com/wealdtech/go-ens"
)

func NewLookupService(clientURL string) nameservice.Lookup {
	client, err := ethclient.Dial(clientURL)
	if err != nil {
		panic(err)
	}
	return &LookupService{
		client: client,
	}
}

type LookupService struct {
	client *ethclient.Client
}

func (s LookupService) ResolveName(ctx context.Context, protocol, network, domainName string) ([]byte, error) {
	address, err := ens.Resolve(s.client, domainName)
	if err != nil {
		return nil, wrapError(err)
	}
	return address.Bytes(), nil
}

func (s LookupService) ResolveAddress(ctx context.Context, protocol, network string, address []byte) (string, error) {
	ethAddress := common.BytesToAddress(address)
	reverse, err := ens.ReverseResolve(s.client, &ethAddress)
	if err != nil {
		return "", wrapError(err)
	}
	return reverse, nil
}
