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

package encoding

// DataPrefix used to identify Mailchain messages.
func DataPrefix() []byte {
	return []byte{0x6d, 0x61, 0x69, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x6e}
}

const (
	// TypeHex encoding value.
	TypeHex = "hex/plain"
	// TypeHex0XPrefix encoding value.
	TypeHex0XPrefix = "hex/0x-prefix"
	// TypeBase58 encoding value.
	TypeBase58 = "base58/plain"
	// TypeBase58SubstrateAddress encoding value.
	TypeBase58SubstrateAddress = "base58/ss58-address"
)
