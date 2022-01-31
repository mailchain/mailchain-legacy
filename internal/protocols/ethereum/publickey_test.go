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

package ethereum

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestGetPublicKeyFromTransaction(t *testing.T) {
	type args struct {
		r        *big.Int
		s        *big.Int
		v        *big.Int
		to       []byte
		input    []byte
		nonce    uint64
		gasPrice *big.Int
		gas      uint64
		value    *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"v=0x29",
			args{
				hexutil.MustDecodeBig("0x3ada323710def1e02f3586710ae3624ceefba1638e9d9894f724a5401997cd79"),
				hexutil.MustDecodeBig("0x2933ddfd0687874e515a8ab479a38646e6db9f3d8b74d27c4e4eae5a116f9f14"),
				hexutil.MustDecodeBig("0x29"),
				hexutil.MustDecode("0xd5ab4ce3605cd590db609b6b5c8901fdb2ef7fe6"),
				hexutil.MustDecode("0x6d61696c636861696e383162336636383539326431393338396439656432346664636338316331666630323835383962653535303436303532366631633961613436623864333739346337653032616565363563386631373733376361366637333564393565303965366131396636303838366638313239326535373835373133343562386531653466393238326531306433396637316238636639653731613231656336393939333637346634616261643231623831393531646565346665643565666465663334643131303264346333336538626662613330623461343730646162643434653938653262363439346136653862363963393336353864393631393639356633313561356266356262313865363265336266623237363463363335323631616366363730303862353761316262333838353164396132656635353730323861336166373839646537396234346662346130336137653637393037343030376531623237"),
				hexutil.MustDecodeUint64("0x8"),
				hexutil.MustDecodeBig("0x12a05f200"),
				hexutil.MustDecodeUint64("0xb274"),
				hexutil.MustDecodeBig("0x30"),
			},
			hexutil.MustDecode("0xbdf6fb97c97c126b492186a4d5b28f34f0671a5aacc974da3bde0be93e45a1c50f89ceff72bd04ac9e25a04a1a6cb010aedaf65f91cec8ebe75901c49b63355d"),
			false,
		},
		{
			"v=0x1c",
			args{
				hexutil.MustDecodeBig("0x6e28ef7db73cd58e9071a411412510402e2090c32a4d81a694d63b67b6ed37a"),
				hexutil.MustDecodeBig("0x411113c3d3f1cadf2b068b224f94ce6fe003fcb9ef9be44b01088767ee8d5cf6"),
				hexutil.MustDecodeBig("0x1c"),
				hexutil.MustDecode("0x92d8f10248c6a3953cc3692a894655ad05d61efb"),
				hexutil.MustDecode("0x"),
				hexutil.MustDecodeUint64("0xc9409"),
				hexutil.MustDecodeBig("0x3b9aca00"),
				hexutil.MustDecodeUint64("0x4cb26"),
				hexutil.MustDecodeBig("0xde0b6b3a7640000"),
			},
			hexutil.MustDecode("0x0bd518dd837e6ed3b902452c0075a4f8d09c8a194cf0ecb8012ca419b6f13916ca560cc840413edcd8cd91c43ca6d86a2d1e8b0bd1bb5fa2c35044fbb42a3cd1"),
			false,
		},
		{
			"hash=0x9220257407f78ad91f340f856fe147751a95257783c0c2c288a129d356ab05e4",
			args{
				hexutil.MustDecodeBig("0x5d04917ff4c0cb832088e1f38bc1b98fd9ebc35ec565fa6475f0b1fdca392aea"),
				hexutil.MustDecodeBig("0x5eb0eac05c23a1cc99af6b1074ef1a97a921f4b50c7b000af27d64630c773922"),
				hexutil.MustDecodeBig("0x26"),
				hexutil.MustDecode("0x5d086f15b2037d2a2be5a0bc2cb2b8472bd0212b"),
				hexutil.MustDecode("0x"),
				hexutil.MustDecodeUint64("0x42138"),
				hexutil.MustDecodeBig("0x165a0bc00"),
				hexutil.MustDecodeUint64("0x186a0"),
				hexutil.MustDecodeBig("0x1908b21b70d456b"),
			},
			hexutil.MustDecode("0xccad0a3df2efc8d4965f317b6f4ad0140e58d21ea2c2d81473b9073a9485aaa4c491a521d6b9fe9d117d926ac278ac56d8415700b239dfe05867305cfab9efa7"),
			false,
		},
		{
			"hash=0x0248c7f152207e9402bf40622f05c7eb24215153417a568c5ecfa1706d68a118",
			args{
				hexutil.MustDecodeBig("0xb4f369be668f83f3be9db16ddcb76baf5215e88fab0d0957913ee8bfc4cdb31a"),
				hexutil.MustDecodeBig("0xe24b814bcf794d338112738538c32797e23395f857183ce7334a149398f5fe"),
				hexutil.MustDecodeBig("0x26"),
				hexutil.MustDecode("0x389f56acec762b7ab765894474d75fcf654b216b"),
				hexutil.MustDecode("0x"),
				hexutil.MustDecodeUint64("0x9"),
				hexutil.MustDecodeBig("0x4a817c800"),
				hexutil.MustDecodeUint64("0x5208"),
				hexutil.MustDecodeBig("0xd02ab486cedc0000"),
			},
			hexutil.MustDecode("0x1e2fc2abffa85f913bdc41dc69f38542bbeb67a515ed88b2b7873b6690e18d4fd62d4dd81c578d0ce7a6da944807ce03ab9d1af0277e724b98a2ddaaf72fe81c"),
			false,
		},
		{
			"hash=0x3f5e9d09a3ff8144682748b2b218ef0cc79de02a19488b097c0d914f12a13189",
			args{
				r:        hexutil.MustDecodeBig("0x6cd8f09351cee443e718467cd1b0b088df6954101faa431a061b350e06e900d8"),
				s:        hexutil.MustDecodeBig("0x1eb2c28fae56382033e2a6be7a1cff80bba6e5ef2c4fcb43ef69fe17bb3ab6c"),
				v:        hexutil.MustDecodeBig("0x25"),
				to:       hexutil.MustDecode("0x0d8775f648430679a709e98d2b0cb6250d2887ef"),
				input:    hexutil.MustDecode("0xa9059cbb000000000000000000000000705d68a74d35f9d9f1d97362fa562c76f649efa000000000000000000000000000000000000000000000000a4cc799563c380000"),
				nonce:    hexutil.MustDecodeUint64("0x0"),
				gasPrice: hexutil.MustDecodeBig("0x165a0bc00"),
				gas:      hexutil.MustDecodeUint64("0x30d40"),
				value:    hexutil.MustDecodeBig("0x0"),
			},
			hexutil.MustDecode("0x0076e16e31760a85d17442594d2f5be77c8b3aa62952b14fbcad71eb84ae15f0d31ae6235fbed75b6530a49e26227e23a7001b977db0656356a85b83ea134871"),
			false,
		},
		{
			"err-sig-to-pub",
			args{
				hexutil.MustDecodeBig("0x0"),
				hexutil.MustDecodeBig("0x0"),
				hexutil.MustDecodeBig("0x0"),
				hexutil.MustDecode("0x5d086f15b2037d2a2be5a0bc2cb2b8472bd0212b"),
				hexutil.MustDecode("0x"),
				hexutil.MustDecodeUint64("0x42138"),
				hexutil.MustDecodeBig("0x165a0bc00"),
				hexutil.MustDecodeUint64("0x186a0"),
				hexutil.MustDecodeBig("0x1908b21b70d456b"),
			},
			[]byte(nil),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPublicKeyFromTransaction(tt.args.r, tt.args.s, tt.args.v, tt.args.to, tt.args.input, tt.args.nonce, tt.args.gasPrice, tt.args.gas, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublicKeyFromTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("GetPublicKeyFromTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeriveChainID(t *testing.T) {
	cases := []struct {
		name     string
		chainID  *big.Int
		expected *big.Int
	}{
		{"0x1c", hexutil.MustDecodeBig("0x1c"), big.NewInt(0)},
		{"0x29", hexutil.MustDecodeBig("0x29"), big.NewInt(3)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualValues(t, tc.expected, deriveChainID(tc.chainID))
		})
	}
}

func Test_createSignatureToUseInRecovery(t *testing.T) {
	type args struct {
		r *big.Int
		s *big.Int
		v *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantLen int
	}{
		{
			"success",
			args{
				hexutil.MustDecodeBig("0x3ada323710def1e02f3586710ae3624ceefba1638e9d9894f724a5401997cd79"),
				hexutil.MustDecodeBig("0x2933ddfd0687874e515a8ab479a38646e6db9f3d8b74d27c4e4eae5a116f9f14"),
				hexutil.MustDecodeBig("0x29"),
			},
			hexutil.MustDecode("0x3ada323710def1e02f3586710ae3624ceefba1638e9d9894f724a5401997cd792933ddfd0687874e515a8ab479a38646e6db9f3d8b74d27c4e4eae5a116f9f1400"),
			65,
		},
		{
			"0x0248c7f152207e9402bf40622f05c7eb24215153417a568c5ecfa1706d68a118",
			args{
				hexutil.MustDecodeBig("0xb4f369be668f83f3be9db16ddcb76baf5215e88fab0d0957913ee8bfc4cdb31a"),
				hexutil.MustDecodeBig("0xe24b814bcf794d338112738538c32797e23395f857183ce7334a149398f5fe"),
				hexutil.MustDecodeBig("0x26"),
			},
			hexutil.MustDecode("0xb4f369be668f83f3be9db16ddcb76baf5215e88fab0d0957913ee8bfc4cdb31a00e24b814bcf794d338112738538c32797e23395f857183ce7334a149398f5fe01"),
			65,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createSignatureToUseInRecovery(tt.args.r, tt.args.s, tt.args.v)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("createSignatureToUseInRecovery() = %v, want %v", got, tt.want)
			}
			assert.Len(t, got, tt.wantLen)
		})
	}
}
