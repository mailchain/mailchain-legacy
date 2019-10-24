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

package nameservice

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	noResolverErrorMsg           = "no resolver"
	noResolutionErrorMsg         = "No resolution"
	unregisteredNameErrorMsg     = "unregistered name"
	couldNotParseAddressErrorMsg = "could not parse address"
)

// RFC 1035 error interpretation
var (
	ErrFormat   = errors.New("Format Error")
	ErrServFail = errors.New("Server Failure")
	ErrNXDomain = errors.New("Non-Existent Domain")
	ErrNotImp   = errors.New("Not Implemented")
	ErrRefused  = errors.New("Query Refused")
)

var Rfc1035StatusMap = map[error]int{
	ErrFormat:   1,
	ErrServFail: 2,
	ErrNXDomain: 3,
	ErrNotImp:   4,
	ErrRefused:  5,
}

func IsRfc1035Error(err error) bool {
	_, ok := Rfc1035StatusMap[err]
	return ok
}

func WrapError(err error) error {
	if err == nil {
		return nil
	}

	if isErrorOfAnyType(err, []string{noResolverErrorMsg, noResolutionErrorMsg, unregisteredNameErrorMsg}) {
		return ErrNXDomain
	} else if isErrorOfAnyType(err, []string{couldNotParseAddressErrorMsg}) {
		// address related to not being able to part ens address not ethereum address
		return ErrFormat
	}
	return err
}

func isErrorOfAnyType(err error, errorStrings []string) bool {
	for _, errorMsg := range errorStrings {
		if strings.Contains(err.Error(), errorMsg) {
			return true
		}
	}
	return false
}
