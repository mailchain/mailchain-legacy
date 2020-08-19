package substraterpc

import (
	gsrpc "github.com/mailchain/go-substrate-rpc-client"
)

func New(address string) (*SubstrateRPC, error) {
	api, err := gsrpc.NewSubstrateAPI(address)
	if err != nil {
		return &SubstrateRPC{address: address}, nil
	}

	client := SubstrateClient{api: api}

	return &SubstrateRPC{client: client, address: address}, nil
}

type SubstrateRPC struct {
	client  Client
	address string
}
