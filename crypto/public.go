// Copyright 2019 Finobo
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

//go:generate mockgen -source=public.go -package=cryptotest -destination=./cryptotest/public_mock.go
package crypto

// PublicKey definition usable in all mailchain crypto operations
type PublicKey interface {
	// Bytes returns the raw bytes representation of the public key.
	//
	// The returned bytes are used for encrypting, verifying a signature, and locating an address.
	Bytes() []byte
	// Verify verifies whether sig is a valid signature of message.
	Verify(message, sig []byte) bool
}

type ExtendedPublicKey interface {
	Bytes() []byte
	PublicKey() PublicKey
	Derive(index uint32) (ExtendedPublicKey, error)
}
