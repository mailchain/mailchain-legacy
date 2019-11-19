package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/pkg/errors"
)

func publicKeyFinders(s values.Store) *PublicKeyFinders {
	return &PublicKeyFinders{
		clients: map[string]PublicKeyFinderClient{
			defaults.ClientEtherscanNoAuth:    etherscanPublicKeyFinderNoAuth(s),
			defaults.ClientEtherscan:          etherscanPublicKeyFinder(s),
			defaults.SubstratePublicKeyFinder: substratePublicKeyFinder(s),
		},
	}
}

// PublicKeyFinders configuration element.
type PublicKeyFinders struct {
	clients map[string]PublicKeyFinderClient
}

// Produce `mailbox.PublicKeyFinder` based on configuration settings.
func (s PublicKeyFinders) Produce(client string) (mailbox.PubKeyFinder, error) {
	if client == "" {
		return nil, nil
	}
	m, ok := s.clients[client]
	if !ok {
		return nil, errors.Errorf("%s not a supported public key finder", client)
	}
	return m.Produce()
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s PublicKeyFinders) Output() output.Element {
	elements := []output.Element{}
	for _, c := range s.clients {
		elements = append(elements, c.Output())
	}
	return output.Element{
		FullName: "public-key-finders",
		Elements: elements,
	}
}
