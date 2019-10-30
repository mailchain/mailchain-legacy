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

// import (
// 	"fmt"

// 	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/defaults"
// 	"github.com/mailchain/mailchain/cmd/mailchain/internal/setup"
// 	"github.com/mailchain/mailchain/internal/protocols"
// 	"github.com/mailchain/mailchain/internal/protocols/ethereum"
// 	"github.com/spf13/cobra"
// )

// func configChainEthereum(receiverSelector, senderSelector, pubKeyFinderSelector setup.ChainNetworkExistingSelector) *cobra.Command {
// 	chain := chains.Ethereum
// 	cmd := &cobra.Command{
// 		Use:   "ethereum",
// 		Short: "setup ethereum",
// 	}
// 	for _, network := range ethereum.Networks() {
// 		cmd.AddCommand(configChainEthereumNetwork(chain, network, receiverSelector, senderSelector, pubKeyFinderSelector))
// 	}
// 	return cmd
// }

// func configChainEthereumNetwork(chain, network string,
// 	receiverSelector, senderSelector, pubKeyFinderSelector setup.ChainNetworkExistingSelector) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   network,
// 		Short: fmt.Sprintf("setup %s", network),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			receiver, _ := cmd.Flags().GetString("receiver")
// 			pkf, _ := cmd.Flags().GetString("public-key-finder")
// 			sender, _ := cmd.Flags().GetString("sender")
// 			if _, err := receiverSelector.Select(chain, network, receiver); err != nil {
// 				return err
// 			}
// 			if _, err := senderSelector.Select(chain, network, sender); err != nil {
// 				return err
// 			}
// 			if _, err := pubKeyFinderSelector.Select(chain, network, pkf); err != nil {
// 				return err
// 			}

// 			cmd.Printf("%s %s configured using:\n", chain, network)
// 			cmd.Printf("- %s: messages sent from key owner\n", sender)
// 			cmd.Printf("- %s: messages sent to key owner\n", receiver)
// 			cmd.Printf("- %s: looking up addresses\n", pkf)
// 			return nil
// 		},
// 	}
// 	cmd.Flags().String("sender", defaults.EthereumSender, "sender to use for messages sent from key owner")
// 	cmd.Flags().String("receiver", defaults.EthereumReceiver, "receiver to use for messages sent to key owner")
// 	cmd.Flags().String("public-key-finder", defaults.EthereumPublicKeyFinder, "public key finder to use for looking up addresses")

// 	return cmd
// }
