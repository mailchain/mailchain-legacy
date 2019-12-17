package pq

import (
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
)

var protocolUint8 = map[string]uint8{ //nolint:gochecknoglobals
	protocols.Ethereum: 1,
}

var uint8Protocol = map[uint8]string{ //nolint:gochecknoglobals
	1: protocols.Ethereum,
}

var protocolNetworkUint8 = map[string]map[string]uint8{ //nolint:gochecknoglobals
	protocols.Ethereum: {
		ethereum.Mainnet: 1,
		ethereum.Goerli:  2,
		ethereum.Kovan:   3,
		ethereum.Rinkeby: 4,
		ethereum.Ropsten: 5,
	},
}

var uint8ProtocolNetwork = map[string]map[uint8]string{ //nolint:gochecknoglobals
	protocols.Ethereum: {
		1: ethereum.Mainnet,
		2: ethereum.Goerli,
		3: ethereum.Kovan,
		4: ethereum.Rinkeby,
		5: ethereum.Ropsten,
	},
}
