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

	"github.com/mailchain/mailchain/cmd/mailchain/config"
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/internal/pkg/cmd/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" // nolint: depguard
)

func PublicKeyFinder(cmd *cobra.Command, chain, network, pkFinder string) (string, error) {
	pkFinder, err := selectPublicKeyFinder(chain, network, pkFinder)
	if err != nil {
		return "", err
	}
	if err := config.SetPubKeyFinder(chain, network, pkFinder); err != nil {
		return "", err
	}
	return pkFinder, nil
}

func selectPublicKeyFinder(chain, network, pkFinder string) (string, error) {
	if pkFinder != names.Empty {
		return pkFinder, nil
	}
	pkFinder, skipped, err := prompts.SelectItemSkipable(
		"Public Key Finder",
		[]string{names.Etherscan},
		viper.GetString(fmt.Sprintf("chains.%s.networks.%s.pubkey-finder", chain, network)) != "")
	if err != nil || skipped {
		return "", err
	}
	return pkFinder, nil
}
