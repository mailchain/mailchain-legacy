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

package kdf

const (
	// Scrypt name.
	Scrypt string = "scrypt"
)

// StorageKeyDeriver calculates a storage key based on the derivation method
type StorageKeyDeriver interface {
	DeriveKey(options DeriveOpts) ([]byte, error)
}

// DeriveOpts options for deriving a storage key
type DeriveOpts interface {
	KDF() string
}

// DeriveOptionsBuilder used to build DeriveOpts
type DeriveOptionsBuilder func(*DeriveOpts)
