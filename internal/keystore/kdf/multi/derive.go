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

package multi

import (
	"github.com/mailchain/mailchain/internal/keystore/kdf/scrypt"
	"github.com/pkg/errors"
)

// DeriveKey from the options provided create a storage key to securely store protocol private keys.
func DeriveKey(options OptionsBuilders) (storageKey []byte, kdf string, err error) {
	if options.Scrypt != nil {
		storageKey, err = scrypt.DeriveKey(options.Scrypt)
		return storageKey, "scrypt", err
	}
	return nil, "unknown", errors.Errorf("unknown `kdf`")
}
