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
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/names"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/s3store"
	"github.com/pkg/errors"
)

// GetSentStore create all the clients based on configuration
func (s SentStore) Get() (stores.Sent, error) {
	switch s.viper.GetString("storage.sent") {
	case names.S3:
		return s3store.NewSent(
			s.viper.GetString("stores.s3.region"),
			s.viper.GetString("stores.s3.bucket"),
			s.viper.GetString("stores.s3.access-key-id"),
			s.viper.GetString("stores.s3.secret-access-key"),
		)
	case names.Mailchain, "":
		return stores.NewSentStore(), nil
	default:
		return nil, errors.Errorf("unsupported storage client")
	}
}

func (s SentStore) Set(sentType string) error {
	s.viper.Set("storage.sent", sentType)
	switch sentType {
	case names.S3:
		return s.setS3()
	case names.Mailchain:
		return nil
	default:
		return errors.Errorf("unsupported sender store type")
	}
}

func (s SentStore) setS3() error {
	bucket, err := s.requiredInput("bucket")
	if err != nil {
		return err
	}
	region, err := s.requiredInput("region")
	if err != nil {
		return err
	}
	accessKeyID, err := s.requiredInput("access-key-id")
	if err != nil {
		return err
	}
	secretAccessKey, err := s.requiredInput("secret-access-key")
	if err != nil {
		return err
	}

	s.viper.Set("stores.s3.access-key-id", accessKeyID)
	s.viper.Set("stores.s3.secret-access-key", secretAccessKey)
	s.viper.Set("stores.s3.bucket", bucket)
	s.viper.Set("stores.s3.region", region)
	return nil
}
