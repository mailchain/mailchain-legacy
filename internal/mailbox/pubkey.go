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

package mailbox

//go:generate mockgen -source=pubkey.go -package=mailboxtest -destination=./mailboxtest/pubkey_mock.go
import (
	"context"
)

// PubKeyFinder find public key to encrypt message with
type PubKeyFinder interface {
	PublicKeyFromAddress(ctx context.Context, protocol, network string, address []byte) ([]byte, error)
}
