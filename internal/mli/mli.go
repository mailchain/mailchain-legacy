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

package mli

const (
	// Empty identifier for empty Message Location Identifier
	Empty uint64 = 0
	// Mailchain identifier for Mailchain Message Location Identifier
	Mailchain uint64 = 1
	// MLIIPFS identifier for IPFS Message Location Identifier
	IPFS uint64 = 2
)

// ToAddress maps code to a location.
func ToAddress() map[uint64]string {
	return map[uint64]string{
		Mailchain: addrMailchain,
		IPFS:      addrIPFS,
	}
}

const ( // These should not have a trailing slash
	addrMailchain = "https://mcx.mx"
	addrIPFS      = "https://ipfs.io/ipfs"
)
