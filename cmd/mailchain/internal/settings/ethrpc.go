package settings

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/ethrpc2"
)

func ethereumRPC2Sender(s values.Store, network string) *EthereumRPC2 {
	return &EthereumRPC2{
		network: network,
		Address: values.NewDefaultString(fmt.Sprintf("https://relay.mailchain.xyz/json-rpc/ethereum/%s", network),
			s,
			fmt.Sprintf("senders.ethereum-rpc2-%s.address", network),
		),
	}
}

type EthereumRPC2 struct {
	Address values.String
	network string
}

func (e EthereumRPC2) Produce() (sender.Message, error) {
	return ethrpc2.New(e.Address.Get())
}

func (e EthereumRPC2) Supports() map[string]bool {
	return map[string]bool{
		"ethereum/" + e.network: true,
	}
}

func (e EthereumRPC2) Output() output.Element {
	return output.Element{
		FullName: "senders.ethereum-rpc2-" + e.network,
		Attributes: []output.Attribute{
			e.Address.Attribute(),
		},
	}
}
