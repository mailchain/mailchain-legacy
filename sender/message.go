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

package sender

import (
	"context"

	"github.com/mailchain/mailchain/internal/mailbox/signer"
)

//go:generate mockgen -source=message.go -package=sendertest -destination=./sendertest/message_mock.go

// Message is prepared, signed, and sent.
type Message interface {
	Send(ctx context.Context, network string, to []byte, from []byte, data []byte, signer signer.Signer, opts SendOpts) (err error)
}

// SendOpts options for sending a message.
type SendOpts interface{}
