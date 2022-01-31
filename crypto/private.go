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

//go:generate mockgen -source=private.go -package=cryptotest -destination=./cryptotest/private_mock.go
package crypto

// PrivateKey definition usable in all mailchain crypto operations
type PrivateKey interface {
	// Bytes returns the byte representation of the private key
	Bytes() []byte
	// PublicKey from the PrivateKey
	PublicKey() PublicKey
	// Sign signs the message with the key and returns the signature.
	Sign(message []byte) ([]byte, error)
}

type ExtendedPrivateKey interface {
	Bytes() []byte
	PrivateKey() PrivateKey
	Derive(index uint32) (ExtendedPrivateKey, error)
	ExtendedPublicKey() (ExtendedPublicKey, error)
}
