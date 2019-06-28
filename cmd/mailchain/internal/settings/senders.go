package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/chains"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/sender"
	"github.com/pkg/errors"
)

func senders(s values.Store) *Senders {
	return &Senders{
		clients: map[string]SenderClient{
			"ethereum-rpc2-" + ethereum.Goerli:  ethereumRPC2Sender(s, ethereum.Goerli),
			"ethereum-rpc2-" + ethereum.Kovan:   ethereumRPC2Sender(s, ethereum.Kovan),
			"ethereum-rpc2-" + ethereum.Mainnet: ethereumRPC2Sender(s, ethereum.Mainnet),
			"ethereum-rpc2-" + ethereum.Rinkeby: ethereumRPC2Sender(s, ethereum.Rinkeby),
			"ethereum-rpc2-" + ethereum.Ropsten: ethereumRPC2Sender(s, ethereum.Ropsten),
			chains.Ethereum + "-relay":          relaySender(s, chains.Ethereum),
		},
	}
}

type Senders struct {
	clients map[string]SenderClient
}

func (s Senders) Produce(client string) (sender.Message, error) {
	m, ok := s.clients[client]
	if !ok {
		return nil, errors.Errorf("%s not a supported sender", client)
	}
	return m.Produce()
}
