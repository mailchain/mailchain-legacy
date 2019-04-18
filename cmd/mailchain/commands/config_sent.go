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
	"github.com/mailchain/mailchain/cmd/mailchain/config"
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prerun"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func cfgStorageSent() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "sent",
		Short:             "Configure sent storage",
		Long:              `Mailchain stores the sent messages so that the recipient can download them.`,
		PersistentPreRunE: prerun.InitConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.AddCommand(cfgSentStorageS3())
	return cmd
}

func cfgSentStorageS3() *cobra.Command {
	return &cobra.Command{
		Use:      "s3",
		Short:    "Configure s3 sender storage",
		PreRunE:  prerun.InitConfig,
		PostRunE: config.WriteConfig,
		RunE: func(cmd *cobra.Command, args []string) error {
			senderStoreType := names.S3

			if err := config.SetSentStorage(senderStoreType); err != nil {
				return errors.WithStack(err)
			}

			cmd.Printf("Sent store %q configured\n", senderStoreType)
			return nil
		},
	}
}
