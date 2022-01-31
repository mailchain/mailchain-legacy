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

package stores

import (
	"github.com/mailchain/mailchain/internal/mail"
)

//go:generate mockgen -source=state.go -package=statemock -destination=./statemock/state_mock.go

// State stores all the actions that support mailbox functionality
type State interface {
	DeleteMessageRead(messageID mail.ID) error
	PutMessageRead(messageID mail.ID) error
	GetReadStatus(messageID mail.ID) (bool, error)

	PutTransaction(protocol, network string, address []byte, tx Transaction) error
	GetTransactions(protocol, network string, address []byte, skip, limit int32) ([]Transaction, error)
}
