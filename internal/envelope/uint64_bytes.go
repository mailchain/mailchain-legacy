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

package envelope

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

// UInt64Bytes is a new "data type" that combines a variable size integer with a byte array to efficiently store both values.
type UInt64Bytes []byte

// NewUInt64Bytes creates a `UInt64Bytes` with the identifier and data added.
func NewUInt64Bytes(i uint64, data []byte) UInt64Bytes {
	loc := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(loc, i)

	buf := make([]byte, len(loc[:n])+1+len(data))
	buf[0] = byte(n)
	copy(buf[1:], loc[:n])
	copy(buf[1+n:], data)

	return buf
}

// UInt64 extracts the identified portion.
func (u UInt64Bytes) UInt64() (uint64, error) {
	i, _, err := parseUInt64Bytes(u)
	return i, err
}

// Bytes extracts the data portion.
func (u UInt64Bytes) Bytes() ([]byte, error) {
	_, b, err := parseUInt64Bytes(u)
	return b, err
}

// Values returns the identified and data from UInt64Bytes.
// Identified is of type `uint64`.
// Data is of type `[]byte`.
func (u UInt64Bytes) Values() (i uint64, b []byte, err error) {
	return parseUInt64Bytes(u)
}

func parseUInt64Bytes(buf []byte) (i uint64, b []byte, err error) {
	if len(buf) == 0 {
		return 0, []byte{}, errors.Errorf("\"buf\" must not be empty")
	}
	bufLen := int(buf[0])

	if len(buf) < bufLen+1 {
		return 0, []byte{}, errors.Errorf("\"buf\" is too short to be valid")
	}

	intPortion := buf[1 : 1+bufLen]
	v, n := binary.Uvarint(intPortion)

	if n != len(intPortion) {
		return 0, []byte{}, errors.Errorf("uint64 did not consume all data")
	}
	return v, buf[1+bufLen:], nil
}
