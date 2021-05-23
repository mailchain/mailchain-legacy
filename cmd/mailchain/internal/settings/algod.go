package settings

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/algod"
)

func algodSender(s values.Store, network string) *AlgodSender {
	var defaultAddress = map[string]string{
		algorand.Mainnet: "https://api.algoexplorer.io",
		algorand.Testnet: "https://api.testnet.algoexplorer.io",
		algorand.Betanet: "https://api.betanet.algoexplorer.io",
	}

	return &AlgodSender{
		network: network,
		Address: values.NewDefaultString(defaultAddress[network], s, fmt.Sprintf("senders.algod-%s.address", network)),
		Token:   values.NewDefaultString("", s, fmt.Sprintf("senders.algod-%s.token", network)),
	}
}

// AlgodSender configuration element.
type AlgodSender struct {
	Address values.String
	Token   values.String
	network string
}

// Produce `sender.Message` based on configuration settings.
func (e AlgodSender) Produce() (sender.Message, error) {
	return algod.New(e.Address.Get(), e.Token.Get())
}

// Supports a map of what protocol and network combinations are supported.
func (e AlgodSender) Supports() map[string]bool {
	return map[string]bool{
		protocols.Algorand + "/" + e.network: true,
	}
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (e AlgodSender) Output() output.Element {
	return output.Element{
		FullName: "senders.algod-" + e.network,
		Attributes: []output.Attribute{
			e.Address.Attribute(),
		},
	}
}
