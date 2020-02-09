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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseUInt64Bytes(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name      string
		args      args
		wantInt   uint64
		wantBytes []byte
		wantErr   bool
	}{
		{
			"1-byte-number",
			args{
				[]byte{0x1, 0x6, 0x76, 0x61, 0x6c, 0x75, 0x65},
			},
			6,
			[]byte("value"),
			false,
		},
		{
			"2-byte-number",
			args{
				[]byte{0x2, 0xe4, 0x2, 0x76, 0x61, 0x6c, 0x75, 0x65},
			},
			356,
			[]byte("value"),
			false,
		},
		{
			"3-byte-number",
			args{
				[]byte{0x3, 0x94, 0xcb, 0x15, 0x76, 0x61, 0x6c, 0x75, 0x65},
			},
			353684,
			[]byte("value"),
			false,
		},
		{
			"zero-value-int",
			args{
				[]byte{0x1, 0x0, 0x76, 0x61, 0x6c, 0x75, 0x65},
			},
			0,
			[]byte("value"),
			false,
		},
		{
			"no-bytes-data",
			args{
				[]byte{0x2, 0xe4, 0x2},
			},
			356,
			[]byte{},
			false,
		},
		{
			"err-too-short",
			args{
				[]byte{0x2, 0xe4},
			},
			0,
			[]byte{},
			true,
		},
		{
			"err-empty",
			args{
				[]byte{},
			},
			0,
			[]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInt, gotBytes, err := parseUInt64Bytes(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUInt64Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInt != tt.wantInt {
				t.Errorf("parseUInt64Bytes() gotInt = %v, wantInt %v", gotInt, tt.wantInt)
			}
			if !assert.Equal(t, tt.wantBytes, gotBytes) {
				t.Errorf("parseUInt64Bytes() gotBytes = %v, wantBytes %v", gotBytes, tt.wantBytes)
			}
		})
	}
}

func TestNewUInt64Bytes(t *testing.T) {
	type args struct {
		i uint64
		b []byte
	}
	tests := []struct {
		name string
		args args
		want UInt64Bytes
	}{
		{
			"1-byte-int",
			args{
				6,
				[]byte("value"),
			},
			append([]byte{0x1, 0x6}, []byte("value")...),
		},
		{
			"2-byte-int",
			args{
				356,
				[]byte("value"),
			},
			append([]byte{0x2, 0xe4, 0x2}, []byte("value")...),
		},
		{
			"3-byte-int",
			args{
				353684,
				[]byte("value"),
			},
			append([]byte{0x3, 0x94, 0xcb, 0x15}, []byte("value")...),
		},
		{
			"zero-int",
			args{
				0,
				[]byte("value"),
			},
			append([]byte{0x1, 0x0}, []byte("value")...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUInt64Bytes(tt.args.i, tt.args.b)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("NewUInt64Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUInt64Bytes_UInt64(t *testing.T) {
	tests := []struct {
		name    string
		u       UInt64Bytes
		want    uint64
		wantErr bool
	}{
		{
			"success",
			[]byte{0x2, 0xe4, 0x2, 0x76, 0x61, 0x6c, 0x75, 0x65},
			356,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.UInt64()
			if (err != nil) != tt.wantErr {
				t.Errorf("UInt64Bytes.UInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UInt64Bytes.UInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUInt64Bytes_Bytes(t *testing.T) {
	tests := []struct {
		name    string
		u       UInt64Bytes
		want    []byte
		wantErr bool
	}{
		{
			"success",
			[]byte{0x2, 0xe4, 0x2, 0x76, 0x61, 0x6c, 0x75, 0x65},
			[]byte("value"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("UInt64Bytes.Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UInt64Bytes.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUInt64Bytes_Values(t *testing.T) {
	tests := []struct {
		name    string
		u       UInt64Bytes
		wantI   uint64
		wantB   []byte
		wantErr bool
	}{
		{
			"success",
			[]byte{0x2, 0xe4, 0x2, 0x76, 0x61, 0x6c, 0x75, 0x65},
			356,
			[]byte("value"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotI, gotB, err := tt.u.Values()
			if (err != nil) != tt.wantErr {
				t.Errorf("UInt64Bytes.Values() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotI != tt.wantI {
				t.Errorf("UInt64Bytes.Values() gotI = %v, want %v", gotI, tt.wantI)
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("UInt64Bytes.Values() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
