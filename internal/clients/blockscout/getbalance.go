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

package blockscout

import (
	"context"

	"github.com/pkg/errors"
)

// Receive check ethereum transactions for mailchain messages
func (c APIClient) GetBalance(ctx context.Context, protocol, network string, address []byte) (string, error) {
	if !c.isNetworkSupported(network) {
		return "", errors.Errorf("network not supported")
	}

	balance, err := c.getBalanceByAddress(network, address)
	if err != nil {
		return "", errors.WithMessage(err, "could not get balance")
	}

	return balance.balance, nil
}
