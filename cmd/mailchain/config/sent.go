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
	"github.com/mailchain/mailchain/cmd/mailchain/config/names"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/commands/prompts"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/s3"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

// GetSentStorage create all the clients based on configuration
func GetSentStorage() (stores.Sent, error) {
	if viper.GetString("storage.sent") == names.S3 {
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

func SetSentStorage(sentType string) error {
	viper.Set("storage.sent", sentType)
	switch sentType {
	case names.S3:
		return setS3()
	default:
		return errors.Errorf("unsupported sender store type")
	}
}

func setS3() error {
	bucket, err := prompts.RequiredInput("bucket")
	if err != nil {
		return err
	}
	region, err := prompts.RequiredInput("region")
	if err != nil {
		return err
	}
	accessKeyID, err := prompts.RequiredInput("access-key-id")
	if err != nil {
		return err
	}
	secretAccessKey, err := prompts.RequiredInput("secret-access-key")
	if err != nil {
		return err
	}

	viper.Set("stores.s3.access-key-id", accessKeyID)
	viper.Set("stores.s3.secret-access-key", secretAccessKey)
	viper.Set("stores.s3.bucket", bucket)
	viper.Set("stores.s3.region", region)
	return nil
}
