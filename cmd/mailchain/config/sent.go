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

package config

import (
	"github.com/mailchain/mailchain/internal/pkg/stores"
	"github.com/mailchain/mailchain/internal/pkg/stores/s3"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

// GetSenderStorage create all the clients based on configuration
func GetSenderStorage() (stores.Sent, error) {
	if viper.GetString("storage.sent") == "s3" {
		return getS3Client()
	}
	return nil, errors.Errorf("unsupported storage client")
}

func getS3Client() (stores.Sent, error) {
	return s3.NewSentStore(
		viper.GetString("stores.s3.region"),
		viper.GetString("stores.s3.bucket"),
		viper.GetString("stores.s3.access-key-id"),
		viper.GetString("stores.s3.secret-access-key"),
	)
}
