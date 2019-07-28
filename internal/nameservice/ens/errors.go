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

package ens

import (
	"fmt"

	"github.com/mailchain/mailchain/internal/nameservice"
	"github.com/pkg/errors"
)

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	switch fmt.Sprintf("%v", errors.Cause(err)) {
	case fmt.Sprintf("%v", errors.Cause(errors.Errorf("no resolver"))),
		fmt.Sprintf("%v", errors.Cause(errors.Errorf("No resolution"))):
		return errors.WithMessage(err, nameservice.ErrUnableToResolve.Error())

	case fmt.Sprintf("%v", errors.Cause(errors.Errorf("unregistered name"))):
		return errors.WithMessage(err, nameservice.ErrNotFound.Error())

	case fmt.Sprintf("%v", errors.Cause(errors.Errorf("could not parse address"))):
		// address related to not being able to part ens address not ethereum address
		return errors.WithMessage(err, nameservice.ErrInvalidName.Error())

	default:
		return err
	}
}
