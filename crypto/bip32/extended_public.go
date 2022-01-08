package bip32

import (
	"encoding/binary"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
)

type ExtendedPublicKey struct {
	depth             byte
	parentFingerPrint uint32 // [4] bytes
	index             uint32 // also known as child number [4] bytes
	chainCode         [32]byte
	key               secp256k1.PublicKey // [33] bytes
}

func (k *ExtendedPublicKey) Bytes() []byte {
	var childNumBytes, fingerprint [4]byte
	var serializationBytes [serializationKeyLen]byte

	binary.BigEndian.PutUint32(childNumBytes[:], k.index)
	binary.BigEndian.PutUint32(fingerprint[:], k.parentFingerPrint)

	serializationBytes[serializationBeginDepth] = k.depth
	copy(serializationBytes[serializationBeginFingerPrint:serializationEndFingerPrint], fingerprint[:])
	copy(serializationBytes[serializationBeginIndex:serializationEndIndex], childNumBytes[:])
	copy(serializationBytes[serializationBeginChainCode:serializationEndChainCode], k.chainCode[:])
	copy(serializationBytes[serializationBeginKeyBytes:serializationEndKeyBytes], k.key.Bytes())

	return serializationBytes[:]
}

func (k *ExtendedPublicKey) PublicKey() crypto.PublicKey {
	return k.key
}

func (k *ExtendedPublicKey) Derive(index uint32) (crypto.ExtendedPublicKey, error) {
	var fingerprint [4]byte
	binary.BigEndian.PutUint32(fingerprint[:], k.parentFingerPrint)

	child, err := hdkeychain.NewExtendedKey(
		[]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub used for serialization only
		k.key.Bytes(),
		k.chainCode[:],
		fingerprint[:],
		k.depth,
		k.index,
		false,
	).Derive(index)
	if err != nil {
		return nil, err
	}

	return fromExtendedPublicKey(child)
}

func ExtendedPublicKeyFromBytes(in []byte) (*ExtendedPublicKey, error) {
	depth := in[serializationBeginDepth]
	parentFP := in[serializationBeginFingerPrint:serializationEndFingerPrint]
	childNum := binary.BigEndian.Uint32(in[serializationBeginIndex:serializationEndIndex])
	chainCode := in[serializationBeginChainCode:serializationEndChainCode]
	keyData := in[serializationBeginKeyBytes:serializationEndKeyBytes]

	return fromExtendedPublicKey(hdkeychain.NewExtendedKey(versionEmpty[:], keyData, chainCode, parentFP, depth, childNum, false))
}

func fromExtendedPublicKey(in *hdkeychain.ExtendedKey) (*ExtendedPublicKey, error) {
	rawPk, err := in.ECPubKey()
	if err != nil {
		return nil, err
	}

	key, err := secp256k1.PublicKeyFromBytes(rawPk.SerializeCompressed())
	if err != nil {
		return nil, err
	}

	var chainCode [32]byte
	copy(chainCode[:], in.ChainCode())

	return &ExtendedPublicKey{
		key:               *(key.(*secp256k1.PublicKey)),
		chainCode:         chainCode,
		parentFingerPrint: in.ParentFingerprint(),
		index:             in.ChildIndex(),
		depth:             in.Depth(),
	}, nil
}
