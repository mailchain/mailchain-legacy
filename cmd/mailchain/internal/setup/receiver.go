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

	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/names"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" // nolint: depguard
)

func Receiver(cmd *cobra.Command, chain, network, receiver string) (string, error) {
	receiver, err := selectReceiver(chain, network, receiver)
	if err != nil {
		return "", err
	}
	if err := config.DefaultReceiver().Set(chain, network, receiver); err != nil {
		return "", err
	}
	return receiver, nil
}

func selectReceiver(chain, network, receiver string) (string, error) {
	if receiver != names.RequiresValue {
		return receiver, nil
	}
	receiver, skipped, err := prompts.SelectItemSkipable(
		"Receiver",
		[]string{names.EtherscanNoAuth, names.Etherscan},
		viper.GetString(fmt.Sprintf("chains.%s.networks.%s.receiver", chain, network)) != "")
	if err != nil || skipped {
		return "", err
	}
	return receiver, nil
}
