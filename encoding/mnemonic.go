package encoding

import (
	"github.com/algorand/go-algorand/crypto/passphrase"
)

// DecodeMnemonicAlgorand returns the bytes represented by the word list as specified by BIP39 as per the algorand implementation.
//
// DecodeMnemonicAlgorand expects that src contain exactly 25 words, all words MUST be in BIP39 word list and are space delimeted.
// If the input is malformed, DecodeMnemonicAlgorand returns an error.
func DecodeMnemonicAlgorand(src string) ([]byte, error) {
	return passphrase.MnemonicToKey(src)
}

// EncodeMnemonicAlgorand returns the string represented as BIP39 word list as per the Algorand implementation.
//
// EncodeMnemonicAlgorand expects that src contains only base58 byte.
// If the input is malformed, EncodeBase58 returns an error.
func EncodeMnemonicAlgorand(src []byte) (string, error) {
	return passphrase.KeyToMnemonic(src)
}
