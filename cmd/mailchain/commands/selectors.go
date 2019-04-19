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
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/cmd/mailchain/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" // nolint: depguard
)

func selectNetwork(cmd *cobra.Command, args, networks []string) (string, error) {
	flg, _ := cmd.Flags().GetString("network")
	if flg != "" {
		return flg, nil
	}
	if len(args) == 1 {
		return args[0], nil
	}
	return prompts.SelectItem("Network", networks)
}

func selectKeystore() (string, error) {
	keystoreType, skipped, err := prompts.SelectItemSkipable(
		"Key Store",
		[]string{names.KeystoreNACLFilestore},
		viper.GetString("storage.keys") != "")
	if err != nil || skipped {
		return "", err
	}
	fmt.Printf("%s used for storing keys\n", keystoreType)
	return keystoreType, nil
}
