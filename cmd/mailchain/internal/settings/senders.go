package settings

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/mailchain/mailchain/sender"
	"github.com/pkg/errors"
)

func senders(s values.Store) *Senders {
	return &Senders{
		clients: map[string]SenderClient{
			"ethereum-rpc2-" + ethereum.Goerli:             ethereumRPC2Sender(s, ethereum.Goerli),
			"ethereum-rpc2-" + ethereum.Kovan:              ethereumRPC2Sender(s, ethereum.Kovan),
			"ethereum-rpc2-" + ethereum.Mainnet:            ethereumRPC2Sender(s, ethereum.Mainnet),
			"ethereum-rpc2-" + ethereum.Rinkeby:            ethereumRPC2Sender(s, ethereum.Rinkeby),
			"ethereum-rpc2-" + ethereum.Ropsten:            ethereumRPC2Sender(s, ethereum.Ropsten),
			protocols.Ethereum + "-relay":                  relaySender(s, protocols.Ethereum),
			"substrate-rpc-" + substrate.EdgewareDev:       substrateRPCSender(s, substrate.EdgewareDev),
			"substrate-rpc-" + substrate.EdgewareMainnet:   substrateRPCSender(s, substrate.EdgewareMainnet),
			"substrate-rpc-" + substrate.EdgewareBeresheet: substrateRPCSender(s, substrate.EdgewareBeresheet),
		},
	}
}

// Senders configuration element.
type Senders struct {
	clients map[string]SenderClient
}

// Produce `sender.Message` based on configuration settings.
func (s Senders) Produce(client string) (sender.Message, error) {
	if client == "" {
		return nil, nil
	}
	m, ok := s.clients[client]
	if !ok {
		return nil, errors.Errorf("%s not a supported sender", client)
	}
	return m.Produce()
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s Senders) Output() output.Element {
	elements := []output.Element{}
	for _, c := range s.clients {
		elements = append(elements, c.Output())
	}

	return output.Element{
		FullName: "senders",
		Elements: elements,
	}
}
