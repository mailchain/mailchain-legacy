// Copyright 2022 Mailchain Ltd
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

	"github.com/pkg/errors"
)

func CompareContentsToHash(data, hash []byte) error {
	kind, digest, err := parse(hash)
	if err != nil {
		return errors.Wrap(err, "compare hash: parse failed")
	}

	h, err := Create(kind, data)
	if err != nil {
		return errors.Wrap(err, "compare hash: create failed")
	}

	digestContents, err := GetDigest(kind, h)
	if err != nil {
		return errors.Wrap(err, "compare hash: get digest failed")
	}

	if !bytes.Equal(digest, digestContents) {
		return errors.Errorf("compare hash: hashes do not match")
	}

	return nil
}
