// Copyright 2022 Mailchain Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	_ = rlp.Encode(hw, x)
	hw.Sum(h[:0])

	return h
}

// deriveChainID derives the chain id from the given v parameter
func deriveChainID(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}

		return new(big.Int).SetUint64((v - 35) / 2)
	}

	v = new(big.Int).Sub(v, big.NewInt(35))

	return v.Div(v, big.NewInt(2))
}

func prependEmptyBytes(in []byte) []byte {
	var out [32]byte

	copy(out[32-len(in):], in)

	return out[:]
}

func createSignatureToUseInRecovery(r, s, v *big.Int) []byte {
	sig := append(prependEmptyBytes(r.Bytes()), prependEmptyBytes(s.Bytes())...)
	newV := v.Int64()
	chainID := deriveChainID(v)

	if chainID.Int64() > 0 {
		newV -= chainID.Int64()*2 + 8
	}

	recovery := big.NewInt(newV - 27)
	if recovery.Int64() == 0 {
		b := make([]byte, 1)
		sig = append(sig, b...)
	} else {
		sig = append(sig, big.NewInt(recovery.Int64()).Bytes()...)
	}

	return sig
}

func createItems(chainID *big.Int, to, input []byte, nonce uint64, gasPrice *big.Int, gas uint64, value *big.Int) []interface{} {
	if chainID.Uint64() > 0 {
		return []interface{}{
			nonce,
			gasPrice,
			gas,
			to,
			value,
			input,
			chainID.Bytes(),
			"",
			"",
		}
	}

	return []interface{}{
		nonce,
		gasPrice,
		gas,
		to,
		value,
		input,
	}
}

// GetPublicKeyFromTransaction retrieve the public key from the transaction information
func GetPublicKeyFromTransaction(r, s, v *big.Int, to, input []byte, nonce uint64, gasPrice *big.Int, gas uint64, value *big.Int) ([]byte, error) {
	chainID := deriveChainID(v)
	items := createItems(chainID, to, input, nonce, gasPrice, gas, value)
	sig := createSignatureToUseInRecovery(r, s, v)
	hash := rlpHash(items).Bytes()

	recoveredKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return nil, errors.WithMessage(err, "could not convert signature to public key")
	}

	return append(prependEmptyBytes(recoveredKey.X.Bytes()), prependEmptyBytes(recoveredKey.Y.Bytes())...), nil
}
