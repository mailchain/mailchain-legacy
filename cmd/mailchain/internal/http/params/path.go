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

package params

import (
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

func PathMessageID(r *http.Request) (mail.ID, error) {
	id, err := mail.FromHexString(mux.Vars(r)["message_id"])
	if err != nil {
		return nil, err
	}
	if len(id) == 0 {
		return nil, errors.Errorf("must not be empty")
	}
	return id, nil
}

func PathNetwork(r *http.Request) string {
	return strings.ToLower(mux.Vars(r)["network"])
}

func PathAddress(r *http.Request) ([]byte, error) {
	addr := strings.ToLower(mux.Vars(r)["address"])
	if addr == "" {
		return nil, errors.Errorf("'address' must not be empty")
	}
	// TODO: should validate address
	// if !ethereum.IsAddressValid(addr) {
	// 	return nil, "", errors.Errorf("'address' is invalid")
	// }
	// TODO: generic address parsing
	return common.HexToAddress(addr).Bytes(), nil
}
