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

package nameservice

import (
	errs "errors"
	"fmt"
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

// ErrorToRFC1035Status mapping of RFC1035 errors.
func ErrorToRFC1035Status(err error) int {
	switch err {
	case ErrFormat:
		return 1
	case ErrServFail:
		return 2
	case ErrNXDomain:
		return 3
	case ErrNotImp:
		return 4
	case ErrRefused:
		return 5
	case nil:
		return 0
	default:
		return -1
	}
}

func RFC1035StatusToError(status int) error {
	switch status {
	case 0:
		return nil
	case 1:
		return ErrFormat
	case 2:
		return ErrServFail
	case 3:
		return ErrNXDomain
	case 4:
		return ErrNotImp
	case 5:
		return ErrRefused
	default:
		return errs.New("unknown RFC1035 status: " + fmt.Sprintf("%d", status))
	}
}

// WrapError with a RFC1035 error.
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
