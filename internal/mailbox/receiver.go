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

import (
	"context"

	"github.com/mailchain/mailchain/crypto/cipher"
)

type Transaction struct {
	Data    cipher.EncryptedContent
	BlockID []byte
	Hash    []byte
}

//go:generate mockgen -source=receiver.go -package=mailboxtest -destination=./mailboxtest/receiver_mock.go
// Receiver gets encrypted data from blockchain.
type Receiver interface {
	Receive(ctx context.Context, network string, address []byte) ([]Transaction, error)
}

// type ReceiverOpts interface{}
