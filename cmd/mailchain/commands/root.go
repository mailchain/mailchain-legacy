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
	log "github.com/sirupsen/logrus" // nolint: depguard
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mailchain",
		Short: "MailChain node.",
		Long: `Decentralized mailchain client, run it locally.
		Complete documentation is available at xxxx`,
	}
}

// Execute run the command
func Execute() {
	if err := rootCmd().Execute(); err != nil {
		log.Fatalln(err)
	}
}
