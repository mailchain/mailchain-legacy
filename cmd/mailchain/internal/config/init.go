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
	"os"
	"os/user"
	"strings"

	"github.com/mailchain/mailchain/cmd/mailchain/internal/config/defaults"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus" // nolint: depguard
	"github.com/spf13/cobra"
	"github.com/spf13/viper" // nolint: depguard
	"github.com/ttacon/chalk"
)

// MailchainHome set home directory for mailchain
func MailchainHome() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	p := usr.HomeDir + "/.mailchain"

	if _, err := os.Stat(p); os.IsNotExist(err) {
		_ = os.Mkdir(p, 0700)
	}
	return p
}

// Init reads in config file and ENV variables if set.
func Init(viper *viper.Viper, cfgFile, logLevel string) error {
	lvl, err := log.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Warningf("Invalid 'log-level' %q, default to [Warning]", logLevel)
		lvl = log.WarnLevel
	}
	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	viper.SetConfigType(defaults.ConfigFileType)
	viper.SetEnvPrefix("mc")
	viper.SetConfigName(defaults.ConfigFileName)   // name of config file (without extension)
	viper.AddConfigPath(defaults.ConfigPathFirst)  // adding current directory as first search path
	viper.AddConfigPath(defaults.ConfigPathSecond) // adding home directory as second search path
	viper.AutomaticEnv()                           // read in environment variables that match
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func WriteConfig(viper *viper.Viper) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := viper.WriteConfig(); err != nil {
			return errors.WithStack(err)
		}
		cmd.Printf(chalk.Green.Color("Config saved\n"))
		return nil
	}
}
