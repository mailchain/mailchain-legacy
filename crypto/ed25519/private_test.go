package ed25519

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
	"testing/iotest"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKey_Bytes(t *testing.T) {
	tests := []struct {
		name string
		pk   PrivateKey
		want []byte
	}{
		{
			"alice",
			alicePrivateKey,
			alicePrivateKeyBytes,
		},
		{
			"bob",
			bobPrivateKey,
			bobPrivateKeyBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.Bytes(); !assert.Equal(t, tt.want, got) {
				t.Errorf("PrivateKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKeyFromBytes(t *testing.T) {
	type args struct {
		pk []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *PrivateKey
		wantErr bool
	}{
		{
			"success-alice-seed",
			args{
				aliceSeed,
			},
			&alicePrivateKey,
			false,
		},
		{
			"success-alice-bytes",
			args{
				alicePrivateKeyBytes,
			},
			&alicePrivateKey,
			false,
		},
		{
			"success-bob-seed",
			args{
				bobSeed,
			},
			&bobPrivateKey,
			false,
		},
		{
			"success-bob-bytes",
			args{
				bobPrivateKeyBytes,
			},
			&bobPrivateKey,
			false,
		},
		{
			"err-len",
			args{
				encodingtest.MustDecodeHex("39d4c9"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyFromBytes(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("PrivateKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_PublicKey(t *testing.T) {
	tests := []struct {
		name string
		pk   PrivateKey
		want crypto.PublicKey
	}{
		{
			"alice",
			alicePrivateKey,
			&alicePublicKey,
		},
		{
			"bob",
			bobPrivateKey,
			&bobPublicKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.PublicKey(); !assert.Equal(t, tt.want, got) {
				t.Errorf("PrivateKey.PublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_Sign(t *testing.T) {
	tests := []struct {
		name    string
		pk      PrivateKey
		msg     []byte
		want    []byte
		wantErr bool
	}{
		{
			"success-bob",
			bobPrivateKey,
			[]byte("message"),
			[]byte{0x7d, 0x51, 0xea, 0xfa, 0x52, 0x78, 0x31, 0x69, 0xd0, 0xa9, 0x4a, 0xc, 0x9f, 0x2b, 0xca, 0xd5, 0xe0, 0x3d, 0x29, 0x17, 0x33, 0x0, 0x93, 0xf, 0xf3, 0xc7, 0xd6, 0x3b, 0xfd, 0x64, 0x17, 0xae, 0x1b, 0xc8, 0x1f, 0xef, 0x51, 0xba, 0x14, 0x9a, 0xe8, 0xa1, 0xe1, 0xda, 0xe0, 0x5f, 0xdc, 0xa5, 0x7, 0x8b, 0x14, 0xba, 0xc4, 0xcf, 0x26, 0xcc, 0xc6, 0x1, 0x1e, 0x5e, 0xab, 0x77, 0x3, 0xc},
			false,
		},
		{
			"success-alice",
			alicePrivateKey,
			[]byte("egassem"),
			[]byte{0xde, 0x6c, 0x88, 0xe6, 0x9c, 0x9f, 0x93, 0xb, 0x59, 0xdd, 0xf4, 0x80, 0xc2, 0x9a, 0x55, 0x79, 0xec, 0x89, 0x5c, 0xa9, 0x7a, 0x36, 0xf6, 0x69, 0x74, 0xc1, 0xf0, 0x15, 0x5c, 0xc0, 0x66, 0x75, 0x2e, 0xcd, 0x9a, 0x9b, 0x41, 0x35, 0xd2, 0x72, 0x32, 0xe0, 0x54, 0x80, 0xbc, 0x98, 0x58, 0x1, 0xa9, 0xfd, 0xe4, 0x27, 0xc7, 0xef, 0xa5, 0x42, 0x5f, 0xf, 0x46, 0x49, 0xb8, 0xad, 0xbd, 0x5},
			false,
		},
		{
			"err-len",
			PrivateKey{
				key: []byte{0xd, 0x9b, 0x4a, 0x3c, 0x10, 0x72, 0x19, 0x91, 0xc6, 0xb8, 0x6, 0xf0, 0xf3, 0x43, 0x53, 0x5d, 0xc2, 0xb4, 0x6c, 0x74},
			},
			[]byte("message"),
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pk.Sign(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Sign() = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateKey(t *testing.T) {
	type args struct {
		rand io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			args{
				rand.Reader,
			},
			false,
			false,
		},
		{
			"err-rand",
			args{
				iotest.DataErrReader(bytes.NewReader(nil)),
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateKey(tt.args.rand)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("GenerateKey() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}
