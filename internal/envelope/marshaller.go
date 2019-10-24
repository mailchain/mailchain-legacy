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
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

// Unmarshal parses the envelope buffer representation in buf and places the
// decoded result in data.
func Unmarshal(buf []byte) (Data, error) {
	if len(buf) == 0 {
		return nil, errors.Errorf("buffer is empty")
	}
	var err error

	var envData Data

	switch buf[0] {
	case Kind0x01:
		data := &ZeroX01{}
		err = proto.Unmarshal(buf[1:], data)
		envData = data
	case Kind0x50:
		data := &ZeroX50{}
		err = proto.Unmarshal(buf[1:], data)
		envData = data
	default:
		err = errors.Errorf("invalid kind")
	}
	if err != nil {
		return nil, err
	}
	return envData, envData.Valid()
}

// Marshal takes envelope data and encodes it into the wire format,
// returning the data.
func Marshal(data Data) ([]byte, error) {
	switch d := data.(type) {
	case *ZeroX01:
		return prefixedProto(Kind0x01, d)
	default:
		return nil, errors.Errorf("unknown data structure, ")
	}
}

func prefixedProto(kind byte, data proto.Message) ([]byte, error) {
	protoData, err := proto.Marshal(data)
	// send the error later if there was one
	prefixedProto := make([]byte, len(protoData)+1)
	prefixedProto[0] = kind
	copy(prefixedProto[1:], protoData)

	return prefixedProto, errors.WithMessage(err, "could not marshal data")
}
