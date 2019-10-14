package defaults

import (
	"github.com/mailchain/mailchain"
)

type NetworkDefaults struct {
	NameServiceAddress    string
	NameServiceDomainName string
	PublicKeyFinder       string
	Receiver              string
	Sender                string
	Disabled              bool
}

func EthereumNetworkAny() NetworkDefaults {
	return NetworkDefaults{
		NameServiceAddress:    NameServiceAddressKind,
		NameServiceDomainName: NameServiceDomainNameKind,
		PublicKeyFinder:       mailchain.ClientEtherscanNoAuth,
		Receiver:              mailchain.ClientEtherscanNoAuth,
		Sender:                "ethereum-relay",
		Disabled:              false,
	}
}

func SubstrateNetworkAny() NetworkDefaults {
	return NetworkDefaults{
		// NameServiceAddress:    NameServiceAddressKind,
		// NameServiceDomainName: NameServiceDomainNameKind,
		PublicKeyFinder: SubstratePublicKeyFinder,
		// Receiver:              mailchain.ClientEtherscanNoAuth,
		// Sender:                "ethereum-relay",
		Disabled: false,
	}
}
