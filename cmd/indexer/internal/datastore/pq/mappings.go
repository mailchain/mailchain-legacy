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

func getPublicKeyTypeUint8(pub_key_type string) (uint8, error) {
	uPubKeyType, ok := publicKeyTypeUint8[pub_key_type]
	if !ok {
		return 0, errors.Errorf("unknown public_key_type: %q", pub_key_type)
	}
	return uPubKeyType, nil
}

func getPublicKeyTypeString(pub_key_type uint8) (string, error) {
	sPubKeyType, ok := publicKeyTypeString[pub_key_type]
	if !ok {
		return "", errors.Errorf("unknown public_key_type: %d", pub_key_type)
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

var publicKeyTypeUint8 = map[string]uint8{
	crypto.SECP256K1: 1,
	crypto.ED25519:   2,
}

var publicKeyTypeString = make(map[uint8]string)

func init() {
	for key, value := range publicKeyTypeUint8 {
		publicKeyTypeString[value] = key
	}
}
