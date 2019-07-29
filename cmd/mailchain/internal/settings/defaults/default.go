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

package defaults

import (
	"log"
	"path/filepath"

	"github.com/mailchain/mailchain"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const (
	Empty        = ""
	KeystoreKind = "nacl-filestore"

	SentStoreKind             = mailchain.Mailchain
	NameServiceAddressKind    = mailchain.Mailchain
	NameServiceDomainNameKind = mailchain.Mailchain

	MailboxStateKind = "leveldb"

	ConfigFileName   = ".mailchain"
	ConfigSubDirName = ".mailchain"
	ConfigFileKind   = "yaml"

	CORSDisabled = false

	Port = 8080

	// Sender = mailchain.ClientRelay

	// EthereumSender          = Sender
	EthereumReceiver        = mailchain.ClientEtherscanNoAuth
	EthereumPublicKeyFinder = mailchain.ClientEtherscanNoAuth
)

func KeystorePath() string {
	return filepath.Join(MailchainHome(), ".keystore")
}

func MailboxStatePath() string {
	return filepath.Join(MailchainHome(), ".mailbox")
}

func MailchainHome() string {
	d, err := homedir.Dir()
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	return filepath.Join(d, ConfigSubDirName)
}

// // working directory
// dir, err := os.Getwd()
// if err != nil {
// 	return err
// }
