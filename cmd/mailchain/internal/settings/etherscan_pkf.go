// nolint: dupl
package settings

import (
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/mailbox"
)

type EtherscanPublicKeyFinder struct {
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
}

func etherscanPublicKeyFinderNoAuth(s values.Store) *EtherscanPublicKeyFinder {
	return etherscanPublicKeyFinderAny(s, mailchain.ClientEtherscanNoAuth)
}

func etherscanPublicKeyFinder(s values.Store) *EtherscanPublicKeyFinder {
	return etherscanPublicKeyFinderAny(s, mailchain.ClientEtherscan)
}

func etherscanPublicKeyFinderAny(s values.Store, kind string) *EtherscanPublicKeyFinder {
	enabledNetworks := []string{}
	for _, n := range ethereum.Networks() {
		enabledNetworks = append(enabledNetworks, "ethereum/"+n)
	}
	return &EtherscanPublicKeyFinder{
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"public-key-finders."+kind+".enabled-networks",
		),
		APIKey: values.NewDefaultString("", s, "public-key-finders."+kind+".api-key"),
	}
}

func (r EtherscanPublicKeyFinder) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

func (r EtherscanPublicKeyFinder) Produce() (mailbox.PubKeyFinder, error) {
	return etherscan.NewAPIClient(r.APIKey.Get())
}
