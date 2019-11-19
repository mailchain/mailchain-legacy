package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

func receivers(s values.Store) *Receivers {
	return &Receivers{
		clients: map[string]ReceiverClient{
			defaults.ClientEtherscanNoAuth: etherscanReceiverNoAuth(s),
			defaults.ClientEtherscan:       etherscanReceiver(s),
		},
	}
}

// Receivers configuration element.
type Receivers struct {
	clients map[string]ReceiverClient
}

// Produce `mailbox.Receiver` based on configuration settings.
func (s Receivers) Produce(client string) (mailbox.Receiver, error) {
	if client == "" {
		return nil, nil
	}
	m, ok := s.clients[client]
	if !ok {
		return nil, errors.Errorf("%s not a supported receiver", client)
	}
	return m.Produce()
}

// Output configuration as an `output.Element` for use in exporting configuration.
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
