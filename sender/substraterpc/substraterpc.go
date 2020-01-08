package substraterpc

import (
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type Options struct {
	Xt types.Extrinsic
}

//func New(address string) (*SubstrateRPC, error) {
//	api, err := gsrpc.NewSubstrateAPI(address)
//	if err != nil {
//		return nil, err
//	}
//	client := SubstrateClient{api: *api}
//	return &SubstrateRPC{client: client}, nil
//}

type SubstrateRPC struct {
	client Client
}
