package settings

import (
	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

func balanceFinders(s values.Store) *BalanceFinders {
	return &BalanceFinders{
		clients: map[string]BalanceFinderClient{
			defaults.ClientEtherscanNoAuth:    etherscanBalanceFinderNoAuth(s),
			defaults.ClientEtherscan:          etherscanBalanceFinder(s),
			defaults.ClientBlockscoutNoAuth:   blockscoutBalanceFinderNoAuth(s),
			defaults.SubstratePublicKeyFinder: EtherscanBalanceFinder(s),
		},
	}
}

// BalanceFinders configuration element.
type BalanceFinders struct {
	clients map[string]BalanceFinderClient
}

// Produce `mailbox.PublicKeyFinder` based on configuration settings.
func (s BalanceFinders) Produce(client string) (mailbox.Balance, error) {
	if client == "" {
		return nil, nil
	}
	m, ok := s.clients[client]
	if !ok {
		return nil, errors.Errorf("%s not a supported balance finder", client)
	}
	return m.Produce()
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s BalanceFinders) Output() output.Element {
	elements := []output.Element{}
	for _, c := range s.clients {
		elements = append(elements, c.Output())
	}
	return output.Element{
		FullName: "balance-finders",
		Elements: elements,
	}
}
