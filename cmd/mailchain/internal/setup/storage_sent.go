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

package setup

import (
	"github.com/mailchain/mailchain/cmd/mailchain/config"
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" // nolint: depguard
)

func SentStorage(cmd *cobra.Command, sentStorageType string) (string, error) {
	sentStorageType, err := selectSentStorage(sentStorageType)
	if err != nil {
		return "", err
	}
	if err := config.DefaultSentStore().Set(sentStorageType); err != nil {
		return "", err
	}
	return sentStorageType, nil
}

func selectSentStorage(sentStorageType string) (string, error) {
	if sentStorageType != names.Empty {
		return sentStorageType, nil
	}
	sentStorageType, skipped, err := prompts.SelectItemSkipable(
		"Sent Store",
		[]string{names.S3},
		viper.GetString("storage.sent") != "")
	if err != nil || skipped {
		return "", err
	}
	return sentStorageType, nil
}
