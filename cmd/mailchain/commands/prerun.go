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

package commands

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" //nolint: depguard
)

func prerunInitConfig(v *viper.Viper) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfgFile, _ := cmd.Flags().GetString("config")
		logLevel, _ := cmd.Flags().GetString("log-level")
		_ = viper.BindPFlag("logger.level", cmd.Flags().Lookup("log-level"))
		preventInitConfig, _ := cmd.Flags().GetBool("prevent-init-config")

		return settings.InitStore(v, cfgFile, logLevel, !preventInitConfig)
	}
}
