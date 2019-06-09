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
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/names"
	"github.com/spf13/viper" // nolint: depguard
)

type SentStorage struct {
	sentStoreSetter    config.SentStoreSetter
	viper              *viper.Viper
	selectItemSkipable func(label string, items []string, skipable bool) (selected string, skipped bool, err error)
}

func (s SentStorage) Select(sentStorageType string) (string, error) {
	sentStorageType, err := s.selectSentStorage(sentStorageType)
	if err != nil {
		return "", err
	}
	if sentStorageType == "" {
		return "", nil
	}
	if err := s.sentStoreSetter.Set(sentStorageType); err != nil {
		return "", err
	}
	return sentStorageType, nil
}

func (s SentStorage) selectSentStorage(sentStorageType string) (string, error) {
	if sentStorageType != names.RequiresValue {
		return sentStorageType, nil
	}
	sentStorageType, skipped, err := s.selectItemSkipable(
		"Sent Store",
		[]string{names.Mailchain, names.S3},
		s.viper.GetString("storage.sent") != "")
	if err != nil || skipped {
		return "", err
	}
	return sentStorageType, nil
}
