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

package keystore

import (
	"fmt"
	"sort"
)

type ProtocolAddress struct {
	Address  []byte
	Protocol string
	Network  string
}

func FlattenAddressesMap(in map[string]map[string][][]byte) []ProtocolAddress {
	out := []ProtocolAddress{}

	for protocol, networkAddresses := range in {
		for network, addresses := range networkAddresses {
			for _, address := range addresses {
				out = append(out, ProtocolAddress{Address: address, Protocol: protocol, Network: network})
			}
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return fmt.Sprintf("%s.%s.%s", out[i].Protocol, out[i].Network, out[i].Address) <
			fmt.Sprintf("%s.%s.%s", out[j].Protocol, out[j].Network, out[j].Address)
	})

	return out
}
