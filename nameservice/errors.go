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
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrUnableToResolve = errors.New("unable to resolve")
	ErrNotFound        = errors.New("not found")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidAddress  = errors.New("invalid address")
)

func IsNoResolverError(err error) bool {
	message := fmt.Sprintf("%v", err)
	if strings.HasPrefix(message, ErrUnableToResolve.Error()) {
		return true
	}
	return false
}

func IsNotFoundError(err error) bool {
	message := fmt.Sprintf("%v", err)
	if strings.HasPrefix(message, ErrNotFound.Error()) {
		return true
	}

	return false
}

func IsInvalidNameError(err error) bool {
	message := fmt.Sprintf("%v", err)
	if strings.HasPrefix(message, ErrInvalidName.Error()) {
		return true
	}

	return false
}

func IsInvalidAddressError(err error) bool {
	message := fmt.Sprintf("%v", err)
	if strings.HasPrefix(message, ErrInvalidAddress.Error()) {
		return true
	}
	return false
}
