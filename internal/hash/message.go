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

// CreateMessageHash used to verify if the contents of the message match the hash.
func CreateMessageHash(encodedData []byte) multihash.Multihash {
	// No err: SHA3_256 does not error
	hash, _ := Create(SHA3256, encodedData)
	return hash
}
