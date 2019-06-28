package settings

import (
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/internal/chains/ethereum"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/mailbox"
)

type EtherscanReceiver struct {
	EnabledProtocolNetworks values.StringSlice
	APIKey                  values.String
}

func etherscanReceiverNoAuth(s values.Store) *EtherscanReceiver {
	return etherscanReceiverAny(s, mailchain.ClientEtherscanNoAuth)
}

func etherscanReceiver(s values.Store) *EtherscanReceiver {
	return etherscanReceiverAny(s, mailchain.ClientEtherscan)
}

func etherscanReceiverAny(s values.Store, kind string) *EtherscanReceiver {
	enabledNetworks := []string{}
	for _, n := range ethereum.Networks() {
		enabledNetworks = append(enabledNetworks, "ethereum/"+n)
	}
	return &EtherscanReceiver{
		EnabledProtocolNetworks: values.NewDefaultStringSlice(
			enabledNetworks,
			s,
			"receivers."+kind+".enabled-networks",
		),
		APIKey: values.NewDefaultString("", s, "receivers."+kind+".api-key"),
	}
}

func (r EtherscanReceiver) Supports() map[string]bool {
	m := map[string]bool{}
	for _, np := range r.EnabledProtocolNetworks.Get() {
		m[np] = true
	}
	return m
}

func (r EtherscanReceiver) Produce() (mailbox.Receiver, error) {
	return etherscan.NewAPIClient(r.APIKey.Get())
}
