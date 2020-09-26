package settings

import (
	"fmt"

	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/substraterpc"
)

func substrateRPCSender(s values.Store, network string) *SubstrateRPC {
	var defaultAddress = map[string]string{
		substrate.EdgewareDev:       "ws://localhost:9944",
		substrate.EdgewareMainnet:   "ws://mainnet1.edgewa.re:9944",
		substrate.EdgewareBeresheet: "ws://beresheet1.edgewa.re:9944",
	}

	return &SubstrateRPC{
		network: network,
		Address: values.NewDefaultString(defaultAddress[network], s, fmt.Sprintf("senders.substrate-rpc-%s.address", network)),
	}
}

// SubstrateRPC configuration element.
type SubstrateRPC struct {
	Address values.String
	network string
}

// Produce `sender.Message` based on configuration settings.
func (e SubstrateRPC) Produce() (sender.Message, error) {
	return substraterpc.New(e.Address.Get())
}

// Supports a map of what protocol and network combinations are supported.
func (e SubstrateRPC) Supports() map[string]bool {
	return map[string]bool{
		"substrate/" + e.network: true,
	}
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (e SubstrateRPC) Output() output.Element {
	return output.Element{
		FullName: "senders.substrate-rpc-" + e.network,
		Attributes: []output.Attribute{
			e.Address.Attribute(),
		},
	}
}
