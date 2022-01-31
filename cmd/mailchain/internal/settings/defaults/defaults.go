// Copyright 2022 Mailchain Ltd.
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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const (
	// Mailchain value.
	Mailchain = "mailchain"

	// Empty default value.
	Empty = ""
	// KeystoreKind default value.
	KeystoreKind = "nacl-filestore"
	// SentStoreKind default value.
	SentStoreKind = Mailchain
	// NameServiceAddressKind default value.
	NameServiceAddressKind = Mailchain
	// NameServiceDomainNameKind default value.
	NameServiceDomainNameKind = Mailchain
	// MailboxStateKind default value.
	MailboxStateKind = "badgerdb"
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
	// AlgorandPublicKeyFinder default value.
	AlgorandPublicKeyFinder = "algorand-public-key-extractor"
	// SubstratePublicKeyFinder default value.
	SubstratePublicKeyFinder = "substrate-public-key-finder"
	// SubstrateBalanceFinder default value.
	SubstrateBalanceFinder = "substrate-balance-finder"
	// CacheTimeout default value for caching interval in memory
	CacheTimeout = "1h"
)

const (
	// ClientBlockscoutNoAuth blockscout without authentication client name.
	ClientBlockscoutNoAuth = "blockscout-no-auth"
	// ClientEtherscan etherscan client name.
	ClientEtherscan = "etherscan"
	// ClientEtherscanNoAuth etherscan without authentication client name.
	ClientEtherscanNoAuth = "etherscan-no-auth"
	// ClientEthereumRPC2 etherscan JSON RPC 2.0 client name.
	ClientEthereumRPC2 = "ethereum-rpc2"
	// ClientRelay relay client name.
	ClientRelay = "relay"
	// ClientAlgod relay client name.
	ClientAlgod = "algod"

	// EthereumRelay relay for ethereum only
	EthereumRelay = "ethereum-relay"
)

// LogsPath default value.
func LogsPath() string {
	return filepath.Join(MailchainHome(), "logs")
}

// KeystorePath default value.
func KeystorePath() string {
	return filepath.Join(MailchainHome(), ".keystore")
}

// MailboxStatePath default value.
func MailboxStatePath() string {
	return filepath.Join(MailchainHome(), ".mailbox")
}

// MessageCachePath default value path for message caching.
func MessageCachePath() string {
	return filepath.Join(MailchainHome(), ".message-cache")
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
