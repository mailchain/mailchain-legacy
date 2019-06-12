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

// nolint: dupl
package setup

import (
	"fmt"

	"github.com/mailchain/mailchain"
)

func (f PubKeyFinder) Select(chain, network, existingPKFinder string) (string, error) {
	pkFinder, err := f.selectPubKeyFinder(chain, network, existingPKFinder)
	if err != nil {
		return "", err
	}
	if pkFinder == "" {
		return "", nil
	}
	if err := f.setter.Set(chain, network, pkFinder); err != nil {
		return "", err
	}
	return pkFinder, nil
}

func (f PubKeyFinder) selectPubKeyFinder(chain, network, existingPKFinder string) (string, error) {
	if existingPKFinder != mailchain.RequiresValue {
		return existingPKFinder, nil
	}
	pkFinder, skipped, err := f.selectItemSkipable(
		"Public Key Finder",
		[]string{mailchain.ClientEtherscanNoAuth, mailchain.ClientEtherscan},
		f.viper.GetString(fmt.Sprintf("chains.%s.networks.%s.pubkey-finder", chain, network)) != "")
	if err != nil || skipped {
		return "", err
	}
	return pkFinder, nil
}
