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

package encoding

// DataPrefix used to identify Mailchain messages.
func DataPrefix() []byte {
	return []byte{0x6d, 0x61, 0x69, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x6e}
}

const (
	// KindHex encoding value.
	KindHex = "hex/plain"
	// KindHex0XPrefix encoding value.
	KindHex0XPrefix = "hex/0x-prefix"
	// KindBase32 encoding value.
	KindBase32 = "base32/plain"
	// KindBase58 encoding value.
	KindBase58 = "base58/plain"
	// KindBase58SubstrateAddress encoding value.
	KindBase58SubstrateAddress = "base58/ss58-address"
	// KindMnemonicAlgorand encoding value.
	KindMnemonicAlgorand = "mnemonic/algorand"
)
