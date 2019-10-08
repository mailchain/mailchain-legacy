package settings

import (
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

func receivers(s values.Store) *Receivers {
	return &Receivers{
		clients: map[string]ReceiverClient{
			mailchain.ClientEtherscanNoAuth: etherscanReceiverNoAuth(s),
			mailchain.ClientEtherscan:       etherscanReceiver(s),
		},
	}
}

type Receivers struct {
	clients map[string]ReceiverClient
}

func (s Receivers) Produce(client string) (mailbox.Receiver, error) {
	m, ok := s.clients[client]
	if !ok {
		return nil, errors.Errorf("%s not a supported receiver", client)
	}
	return m.Produce()
}

func (s Receivers) Output() output.Element {
	elements := []output.Element{}
	for _, c := range s.clients {
		elements = append(elements, c.Output())
	}

	return output.Element{
		FullName: "receivers",
		Elements: elements,
	}
}
