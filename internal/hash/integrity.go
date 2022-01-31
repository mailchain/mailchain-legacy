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

import "github.com/multiformats/go-multihash"

// CreateIntegrityHash returns a hash of the encrypted `[]byte` to allow easy checking it has not been tampered with.
func CreateIntegrityHash(encryptedData []byte) multihash.Multihash {
	hash, _ := Create(MurMur3128, encryptedData)
	return hash
}
