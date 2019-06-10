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

func Sender(cmd *cobra.Command, chain, network, sender string) (string, error) {
	sender, err := selectSender(chain, network, sender)
	if err != nil {
		return "", err
	}
	if err := config.DefaultSender().Set(chain, network, sender); err != nil {
		return "", err
	}
	return sender, nil
}

func selectSender(chain, network, sender string) (string, error) {
	if sender != names.RequiresValue {
		return sender, nil
	}
	sender, skipped, err := prompts.SelectItemSkipable(
		"Sender",
		[]string{names.EthereumRPC2},
		viper.GetString(fmt.Sprintf("chains.%s.networks.%s.sender", chain, network)) != "")
	if err != nil || skipped {
		return "", err
	}
	return sender, nil
}
