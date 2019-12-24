package pq

import (
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
)

func getProtocolNetworkUint8(prot, net string) (protocol, network uint8, err error) {
	uProtocol, ok := protocolUint8[prot]
	if !ok {
		return 0, 0, errors.Errorf("unknown protocol: %q", prot)
	}

	uNetwork, ok := protocolNetworkUint8[prot][net]
	if !ok {
		return 0, 0, errors.Errorf("unknown protocol.network: \"%s.%s\"", prot, net)
	}

	return uProtocol, uNetwork, nil
}

func getPublicKeyTypeUint8(pubKeyType string) (uint8, error) {
	uPubKeyType, ok := publicKeyTypeUint8[pubKeyType]
	if !ok {
		return 0, errors.Errorf("unknown public_key_type: %q", pubKeyType)
	}

	return uPubKeyType, nil
}

func getPublicKeyTypeString(pubKeyType uint8) (string, error) {
	sPubKeyType, ok := publicKeyTypeString[pubKeyType]
	if !ok {
		return "", errors.Errorf("unknown public_key_type: %d", pubKeyType)
	}

	return sPubKeyType, nil
}

var protocolUint8 = map[string]uint8{ //nolint:gochecknoglobals
	protocols.Ethereum: 1,
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

var publicKeyTypeUint8 = map[string]uint8{ //nolint:gochecknoglobals
	crypto.SECP256K1: 1,
	crypto.ED25519:   2,
}

var publicKeyTypeString = map[uint8]string{ //nolint:gochecknoglobals
	1: crypto.SECP256K1,
	2: crypto.ED25519,
}
