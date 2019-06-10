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
	"fmt"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/names"
	"github.com/mailchain/mailchain/internal/keystore/nacl"
	"github.com/pkg/errors"
	"github.com/spf13/viper" // nolint: depguard
)

//go:generate mockgen -source=keystore.go -package=configtest -destination=./configtest/keystore_mock.go

type KeystoreSetter interface {
	Set(keystoreType, keystorePath string) error
}

type Keystore struct {
	viper                    *viper.Viper
	requiredInputWithDefault func(label string, defaultValue string) (string, error)
}

// GetKeystore create new keystore from config
func (k Keystore) Get() (*nacl.FileStore, error) {
	if k.viper.GetString("storage.keys") == names.KeystoreNACLFilestore {
		fs := nacl.NewFileStore(k.viper.GetString(fmt.Sprintf("stores.%s.path", names.KeystoreNACLFilestore)))
		return &fs, nil
	}

	return nil, errors.Errorf("unknown keystore type")
}

func (k Keystore) Set(keystoreType, keystorePath string) error {
	var err error
	switch keystoreType {
	case names.KeystoreNACLFilestore:
		// NACL only needs to set the path
		err = k.setKeystorePath(keystoreType, keystorePath)
	default:
		err = errors.Errorf("unsupported key store type")
	}
	if err == nil {
		k.viper.Set("storage.keys", keystoreType)
	}

	return err
}

func (k Keystore) setKeystorePath(keystoreType, keystorePath string) error {
	if keystorePath == "" {
		res, err := k.requiredInputWithDefault("path", defaults.KeystorePath)
		if err != nil {
			return err
		}
		keystorePath = res
	}
	k.viper.Set(fmt.Sprintf("stores.%s.path", keystoreType), keystorePath)
	return nil
}
