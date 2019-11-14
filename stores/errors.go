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

package stores

import (
	"fmt"

	"github.com/pkg/errors"
	ldberr "github.com/syndtr/goleveldb/leveldb/errors"
)

var (
	errNotFound = errors.New("not found")
)

// IsNotFoundError checks if the error supplied is an not found error. It checks against known errors from supported stores, S3 and leveldb.
func IsNotFoundError(err error) bool {
	switch fmt.Sprintf("%v", errors.Cause(err)) {
	case fmt.Sprintf("%v", errors.Cause(errNotFound)),
		fmt.Sprintf("%v", errors.Cause(ldberr.ErrNotFound)):
		return true
	default:
		return false
	}
}
