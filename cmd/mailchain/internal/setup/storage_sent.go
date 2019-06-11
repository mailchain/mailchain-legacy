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
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/names"
)

func (s SentStorage) Select(existingSentStorageType string) (string, error) {
	sentStorageType, err := s.selectSentStorage(existingSentStorageType)
	if err != nil {
		return "", err
	}
	if sentStorageType == "" {
		return "", nil
	}
	if err := s.setter.Set(sentStorageType); err != nil {
		return "", err
	}
	return sentStorageType, nil
}

func (s SentStorage) selectSentStorage(existingSentStorageType string) (string, error) {
	if existingSentStorageType != names.RequiresValue {
		return existingSentStorageType, nil
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
