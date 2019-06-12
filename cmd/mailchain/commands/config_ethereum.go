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

package commands

import (
	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup"
	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/spf13/cobra"
)

func configChainEthereum() *cobra.Command {
	return &cobra.Command{
		Use:   "ethereum",
		Short: "setup ethereum",
		RunE: func(cmd *cobra.Command, args []string) error {
			chain := encoding.Ethereum
			network, err := setup.DefaultNetwork().Select(cmd, args, chain, mailchain.RequiresValue)
			if err != nil {
				return err
			}

			cmd.Printf("%s configured\n", network)
			return nil
		},
	}
}
