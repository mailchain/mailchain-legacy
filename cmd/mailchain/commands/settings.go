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
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings"
	"github.com/spf13/cobra"
)

func settingsCmd(config *settings.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Settings of the mailchain application",
	}
	cmd.AddCommand(settingsViewAll(config))
	return cmd
}

func settingsViewAll(config *settings.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View the current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			commentDefaults, _ := cmd.Flags().GetBool("comment-defaults")
			excludeDefaults, _ := cmd.Flags().GetBool("exclude-defaults")
			config.ToYaml(cmd.OutOrStderr(), 2, commentDefaults, excludeDefaults)
			return nil
		},
	}
	cmd.Flags().Bool("comment-defaults", true, "comment out values if the value is the default")
	cmd.Flags().Bool("exclude-defaults", false, "exclude values if the value is the default")
	return cmd
}
