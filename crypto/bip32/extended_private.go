package bip32

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil/hdkeychain"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
)

const (
	serializationBeginDepth       = 0
	serializationEndDepth         = serializationBeginDepth + serializationLengthDepth
	serializationBeginFingerPrint = serializationEndDepth
	serializationEndFingerPrint   = serializationBeginFingerPrint + serializationLengthFingerPrint
	serializationBeginIndex       = serializationEndFingerPrint
	serializationEndIndex         = serializationEndFingerPrint + serializationLengthIndex
	serializationBeginChainCode   = serializationEndIndex
	serializationEndChainCode     = serializationBeginChainCode + serializationLengthChainCode
	serializationBeginKeyBytes    = serializationEndChainCode
	serializationEndKeyBytes      = serializationBeginKeyBytes + serializationLengthKeyBytes

	serializationLengthDepth       = 1
	serializationLengthFingerPrint = 4
	serializationLengthIndex       = 4
	serializationLengthChainCode   = 32
	serializationLengthKeyBytes    = 33

	serializationKeyLen = serializationLengthDepth + serializationLengthFingerPrint + serializationLengthIndex + serializationLengthChainCode + serializationLengthKeyBytes // 74 bytes
)

const (
	// MinSeedBytes is the minimum number of bytes allowed for a seed to
	// a master node.
	MinSeedBytes = 16 // 128 bits

	// MaxSeedBytes is the maximum number of bytes allowed for a seed to
	// a master node.
	MaxSeedBytes = 64 // 512 bits
)

var (
	versionEmpty = [4]byte{}

	// ErrNotAPrivateKey describes an error in which the provided bytes are
	// not valid for a private key. The key data portion must start with 0x0
	// as secp256k1 private keys are 32 bytes yet the space allocated for keys
	// in extended keys is 33 bytes. Meaning a private key first byte is always
	// 0x0. This error indicates the key is invalid and the user should check if
	// it's a public key.
	ErrNotAPrivateKey = errors.New("private key must start with leading 0x0")

	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf("seed length must be between %d and %d "+
		"bits", MinSeedBytes*8, MaxSeedBytes*8)
)

type ExtendedPrivateKey struct {
	depth             byte
	parentFingerPrint uint32 // [4] bytes
	index             uint32 // also known as child number [4] bytes
	chainCode         [32]byte
	key               secp256k1.PrivateKey // [33] bytes
}

// Bytes the format is the functional representation of the extended key.
// It is based on the BIP32 format, minus version and checksum which are
// the first and last 4 bytes.
// BIP 32 = version + depth + fingerprint+ child num + chain code + key data + checksum
// Functional = depth + fingerprint+ child num + chain code + key data.
// BIP32 keys are created by adding the version and calculating the checksum.
func (k ExtendedPrivateKey) Bytes() []byte {
	var childNumBytes, fingerprint [4]byte
	var serializationBytes [serializationKeyLen]byte

	binary.BigEndian.PutUint32(childNumBytes[:], k.index)
	binary.BigEndian.PutUint32(fingerprint[:], k.parentFingerPrint)

	serializationBytes[serializationBeginDepth] = k.depth
	copy(serializationBytes[serializationBeginFingerPrint:serializationEndFingerPrint], fingerprint[:])
	copy(serializationBytes[serializationBeginIndex:serializationEndIndex], childNumBytes[:])
	copy(serializationBytes[serializationBeginChainCode:serializationEndChainCode], k.chainCode[:])

	serializationBytes[serializationBeginKeyBytes] = 0x0
	b := k.key.Bytes()
	copy(serializationBytes[serializationBeginKeyBytes+1:serializationEndKeyBytes], b)

	return serializationBytes[:]
}

func (k ExtendedPrivateKey) PrivateKey() crypto.PrivateKey {
	return k.key
}

func (k ExtendedPrivateKey) Derive(index uint32) (crypto.ExtendedPrivateKey, error) {
	var fingerprint [4]byte
	binary.BigEndian.PutUint32(fingerprint[:], k.parentFingerPrint)

	child, err := hdkeychain.NewExtendedKey(
		[]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub used for serialization only
		k.key.Bytes(),
		k.chainCode[:],
		fingerprint[:],
		k.depth,
		k.index,
		true,
	).Derive(index)
	if err != nil {
		return nil, err
	}

	return fromExtendedPrivateKey(child)
}

func (k ExtendedPrivateKey) ExtendedPublicKey() (crypto.ExtendedPublicKey, error) {
	var fingerprint [4]byte
	binary.BigEndian.PutUint32(fingerprint[:], k.parentFingerPrint)

	pubKey, ok := k.key.PublicKey().(*secp256k1.PublicKey)
	if !ok || pubKey == nil {
		return nil, errors.New("invalid public key")
	}

	return fromExtendedPublicKey(
		hdkeychain.NewExtendedKey(
			[]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv used for serialization only
			pubKey.Bytes(),
			k.chainCode[:],
			fingerprint[:],
			k.depth,
			k.index,
			false,
		),
	)
}

func ExtendedPrivateKeyFromBytes(in []byte) (*ExtendedPrivateKey, error) {
	if len(in) != serializationKeyLen {
		return nil, errors.New("key length must be 74 bytes")
	}

	// Deserialize each of the payload fields.
	depth := in[serializationBeginDepth]
	parentFP := in[serializationBeginFingerPrint:serializationEndFingerPrint]
	childNum := binary.BigEndian.Uint32(in[serializationBeginIndex:serializationEndIndex])
	chainCode := in[serializationBeginChainCode:serializationEndChainCode]
	keyData := in[serializationBeginKeyBytes:serializationEndKeyBytes]

	// The key data is a private key if it starts with 0x00.
	// Compressed pubkeys either start with 0x02 or 0x03.
	if keyData[0] != 0x00 {
		return nil, ErrNotAPrivateKey
	}

	return fromExtendedPrivateKey(hdkeychain.NewExtendedKey(versionEmpty[:], keyData, chainCode, parentFP, depth, childNum, true))
}

func ExtendedPrivateKeyFromSeed(seed []byte) (*ExtendedPrivateKey, error) {
	// source: github.com/btcsuite/btcutil/hdkeychain/extendedkey.go#NewMaster
	// cloned to avoiding adding in github.com/btcsuite/btcd/chaincfg as a
	// dependency for a parameter field
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	// masterKey is the master key used along with a random seed used to generate
	// the master node in the hierarchical tree.
	var masterKey = []byte("Bitcoin seed")

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "Bitcoin seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	hmac512.Write(seed)
	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master secret key
	//   Ir = master chain code
	secretKey := lr[:len(lr)/2]
	chainCode := lr[len(lr)/2:]

	// Ensure the key in usable.
	secretKeyNum := new(big.Int).SetBytes(secretKey)
	if secretKeyNum.Cmp(ethcrypto.S256().Params().N) >= 0 || secretKeyNum.Sign() == 0 {
		return nil, secp256k1.ErrUnusableSeed
	}

	parentFP := []byte{0x00, 0x00, 0x00, 0x00}

	return fromExtendedPrivateKey(hdkeychain.NewExtendedKey(versionEmpty[:], secretKey, chainCode, parentFP, 0, 0, true))
}

func fromExtendedPrivateKey(in *hdkeychain.ExtendedKey) (*ExtendedPrivateKey, error) {
	rawPk, err := in.ECPrivKey()
	if err != nil {
		return nil, err
	}

	ecdsa := rawPk.ToECDSA()

	var chainCode [32]byte
	copy(chainCode[:], in.ChainCode())

	return &ExtendedPrivateKey{
		key:               secp256k1.PrivateKeyFromECDSA(*ecdsa),
		chainCode:         chainCode,
		parentFingerPrint: in.ParentFingerprint(),
		index:             in.ChildIndex(),
		depth:             in.Depth(),
	}, nil
}
