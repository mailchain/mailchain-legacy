package defaults

// NetworkDefaults are the required default values for a network.
type NetworkDefaults struct {
	NameServiceAddress    string
	NameServiceDomainName string
	PublicKeyFinder       string
	BalanceFinder         string
	Receiver              string
	Sender                string
	Disabled              bool
}

// EthereumNetworkAny default values for any Ethereum network.
func EthereumNetworkAny() *NetworkDefaults {
	return &NetworkDefaults{
		NameServiceAddress:    NameServiceAddressKind,
		NameServiceDomainName: NameServiceDomainNameKind,
		PublicKeyFinder:       ClientEtherscanNoAuth,
		BalanceFinder:         ClientEtherscanNoAuth,
		Receiver:              ClientEtherscanNoAuth,
		Sender:                EthereumRelay,
		Disabled:              false,
	}
}

// SubstrateNetworkAny default values for any Substrate network.
func SubstrateNetworkAny(network string) *NetworkDefaults {
	return &NetworkDefaults{
		// NameServiceAddress:    NameServiceAddressKind,
		// NameServiceDomainName: NameServiceDomainNameKind,
		PublicKeyFinder: SubstratePublicKeyFinder,
		//	BalanceFinder:   SubstrateBalanceFinder,
		Receiver: Mailchain,
		Sender:   "substrate-rpc-" + network,
		Disabled: false,
	}
}
