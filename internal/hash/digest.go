// Copyright 2020 Finobo
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

package hash

import (
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

func GetDigest(kind int, hash []byte) ([]byte, error) {
	switch kind {
	case SHA3256, MurMur3128:
		o, err := multihash.Decode(hash)
		if err != nil {
			return nil, err
		}

		return o.Digest, err
	case CIVv1SHA2256Raw:
		c, err := cid.Cast(hash)
		if err != nil {
			return nil, err
		}

		o, _ := multihash.Decode(c.Hash()) // cast statement tests known error conditions

		return o.Digest, err
	default:
		return nil, errors.Errorf("unknown hash kind")
	}
}
