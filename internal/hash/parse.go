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
	"bytes"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

func parse(hash []byte) (kind int, digest []byte, err error) {
	if len(hash) == 0 {
		return Unknown, nil, errors.Errorf("hash can not be empty")
	}

	switch hash[0] {
	case multihash.SHA3_256:
		o, err := multihash.Decode(hash)
		if err != nil {
			return Unknown, nil, err
		}

		return SHA3256, o.Digest, err
	case multihash.MURMUR3_128:
		o, err := multihash.Decode(hash)
		if err != nil {
			return Unknown, nil, err
		}

		return MurMur3128, o.Digest, err
	}

	c, err := cid.Cast(hash)
	if err != nil {
		return Unknown, nil, err
	}

	if bytes.Equal(c.Prefix().Bytes(), []byte{0x01, 0x55, 0x12, 0x20}) {
		o, _ := multihash.Decode(c.Hash()) // cast prevents errors here

		return CIVv1SHA2256Raw, o.Digest, err
	}

	return Unknown, nil, errors.Errorf("unknown hash kind")
}
