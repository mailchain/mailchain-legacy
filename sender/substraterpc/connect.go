package substraterpc

import gsrpc "github.com/centrifuge/go-substrate-rpc-client"

func (s SubstrateRPC) Connect() error {
	if s.client != nil {
		return nil
	}

	api, err := gsrpc.NewSubstrateAPI(s.address)
	if err != nil {
		return err
	}

	s.client = SubstrateClient{api: api}

	return nil
}
