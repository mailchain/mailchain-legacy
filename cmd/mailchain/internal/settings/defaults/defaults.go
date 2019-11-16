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
	// Empty default value.
	Empty = ""
	// KeystoreKind default value.
	KeystoreKind = "nacl-filestore"
	// SentStoreKind default value.
	SentStoreKind = mailchain.Mailchain
	// NameServiceAddressKind default value.
	NameServiceAddressKind = mailchain.Mailchain
	// NameServiceDomainNameKind default value.
	NameServiceDomainNameKind = mailchain.Mailchain
	// MailboxStateKind default value.
	MailboxStateKind = "leveldb"
	// ConfigFileName default value.
	ConfigFileName = ".mailchain"
	// ConfigSubDirName default value.
	ConfigSubDirName = ".mailchain"
	// ConfigFileKind default value.
	ConfigFileKind = "yaml"
	// CORSDisabled default value.
	CORSDisabled = false
	// Port default value.
	Port = 8080
	// SubstratePublicKeyFinder default value.
	SubstratePublicKeyFinder = "substrate-public-key-finder"
)

// KeystorePath default value.
func KeystorePath() string {
	return filepath.Join(MailchainHome(), ".keystore")
}

// MailboxStatePath default value.
func MailboxStatePath() string {
	return filepath.Join(MailchainHome(), ".mailbox")
}

// MailchainHome directory default value.
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
