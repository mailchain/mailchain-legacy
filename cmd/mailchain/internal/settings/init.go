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

package settings

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mailchain/mailchain"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus" //nolint: depguard
	"github.com/spf13/viper"         //nolint: depguard
)

// InitStore creates and loads the configuration storage.
func InitStore(v *viper.Viper, cfgFile, logLevel string, createFile bool) error {
	if cfgFile == "" {
		cfgFile = filepath.Join(defaults.MailchainHome(), defaults.ConfigFileName+"."+defaults.ConfigFileKind)
	}
	lvl, err := log.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Warningf("Invalid 'log-level' %q, default to [Warning]", logLevel)
		lvl = log.WarnLevel
	}

	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	v.SetEnvPrefix("mc")
	v.AutomaticEnv()

	v.SetConfigFile(cfgFile)

	err = v.ReadInConfig()
	_, ok := err.(viper.ConfigFileNotFoundError)
	if ok || err != nil && strings.Contains(err.Error(), "no such file or directory") {
		if createFile {
			return createEmptyFile(v, cfgFile)
		}
		return errors.WithMessage(err, "config creation disabled")
	}
	return err
}

func createEmptyFile(v *viper.Viper, fileName string) error {
	dir, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	v.Set("version", mailchain.Version)

	return v.WriteConfigAs(fileName)
}
