package bip32

import (
	"testing"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestExtendedPublicKey_Bytes(t *testing.T) {
	type fields struct {
		key               secp256k1.PublicKey
		chainCode         [32]byte
		parentFingerPrint uint32
		index             uint32
		depth             byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
			fields{
				depth:             0x3,
				parentFingerPrint: 0xbef5a2f9,
				index:             0x80000002,
				chainCode:         [32]uint8{0x4, 0x46, 0x6b, 0x9c, 0xc8, 0xe1, 0x61, 0xe9, 0x66, 0x40, 0x9c, 0xa5, 0x29, 0x86, 0xc5, 0x84, 0xf0, 0x7e, 0x9d, 0xc8, 0x1f, 0x73, 0x5d, 0xb6, 0x83, 0xc3, 0xff, 0x6e, 0xc7, 0xb1, 0x50, 0x3f},
				key: func() secp256k1.PublicKey {
					o, _ := secp256k1.PublicKeyFromBytes(encodingtest.MustDecodeHex("0357bfe1e341d01c69fe5654309956cbea516822fba8a601743a012a7896ee8dc2"))
					d := o.(*secp256k1.PublicKey)
					return *d
				}(),
			},
			[]byte{0x3, 0xbe, 0xf5, 0xa2, 0xf9, 0x80, 0x0, 0x0, 0x2, 0x4, 0x46, 0x6b, 0x9c, 0xc8, 0xe1, 0x61, 0xe9, 0x66, 0x40, 0x9c, 0xa5, 0x29, 0x86, 0xc5, 0x84, 0xf0, 0x7e, 0x9d, 0xc8, 0x1f, 0x73, 0x5d, 0xb6, 0x83, 0xc3, 0xff, 0x6e, 0xc7, 0xb1, 0x50, 0x3f, 0x3, 0x57, 0xbf, 0xe1, 0xe3, 0x41, 0xd0, 0x1c, 0x69, 0xfe, 0x56, 0x54, 0x30, 0x99, 0x56, 0xcb, 0xea, 0x51, 0x68, 0x22, 0xfb, 0xa8, 0xa6, 0x1, 0x74, 0x3a, 0x1, 0x2a, 0x78, 0x96, 0xee, 0x8d, 0xc2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &ExtendedPublicKey{
				key:               tt.fields.key,
				chainCode:         tt.fields.chainCode,
				parentFingerPrint: tt.fields.parentFingerPrint,
				index:             tt.fields.index,
				depth:             tt.fields.depth,
			}
			if got := k.Bytes(); !assert.Equal(t, tt.want, got) {
				t.Errorf("ExtendedPublic.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtendedPublicKeyFromBytes(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *ExtendedPublicKey
		wantErr bool
	}{
		{
			"xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
			args{
				encodingtest.MustDecodeBase58("xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5")[4:78],
			},
			&ExtendedPublicKey{
				depth:             0x3,
				parentFingerPrint: 0xbef5a2f9,
				index:             0x80000002,
				chainCode:         [32]uint8{0x4, 0x46, 0x6b, 0x9c, 0xc8, 0xe1, 0x61, 0xe9, 0x66, 0x40, 0x9c, 0xa5, 0x29, 0x86, 0xc5, 0x84, 0xf0, 0x7e, 0x9d, 0xc8, 0x1f, 0x73, 0x5d, 0xb6, 0x83, 0xc3, 0xff, 0x6e, 0xc7, 0xb1, 0x50, 0x3f},
				key: func() secp256k1.PublicKey {
					o, _ := secp256k1.PublicKeyFromBytes(encodingtest.MustDecodeHex("0357bfe1e341d01c69fe5654309956cbea516822fba8a601743a012a7896ee8dc2"))
					d := o.(*secp256k1.PublicKey)
					return *d
				}(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtendedPublicKeyFromBytes(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtendedPublicKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ExtendedPublicKeyFromBytes() = %v, want %v", got, tt.want)
			}

			if !tt.wantErr && got == nil {
				t.Error("object expected nil returned")
				t.FailNow()
			}

			if !assert.Equal(t, tt.want.PublicKey().Bytes(), got.PublicKey().Bytes()) {
				t.Errorf("ExtendedPublicKeyFromBytes().PublicKey().Bytes() = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.chainCode, got.chainCode) {
				t.Errorf("ExtendedPublicKeyFromBytes().chainCode = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.depth, got.depth) {
				t.Errorf("ExtendedPublicKeyFromBytes().depth = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.index, got.index) {
				t.Errorf("ExtendedPublicKeyFromBytes().index = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.parentFingerPrint, got.parentFingerPrint) {
				t.Errorf("ExtendedPublicKeyFromBytes().parentFingerPrint = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_fromExtendedPublicKey(t *testing.T) {
	type args struct {
		in *hdkeychain.ExtendedKey
	}
	tests := []struct {
		name    string
		args    args
		want    *ExtendedPublicKey
		wantErr bool
	}{
		{
			"xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
			args{
				func() *hdkeychain.ExtendedKey {
					o, _ := hdkeychain.NewKeyFromString("xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5")
					return o
				}(),
			},
			&ExtendedPublicKey{
				key: func() secp256k1.PublicKey {
					o, _ := secp256k1.PublicKeyFromBytes(encodingtest.MustDecodeHex("0357bfe1e341d01c69fe5654309956cbea516822fba8a601743a012a7896ee8dc2"))
					d := o.(*secp256k1.PublicKey)
					return *d
				}(),
				chainCode:         [32]byte{4, 70, 107, 156, 200, 225, 97, 233, 102, 64, 156, 165, 41, 134, 197, 132, 240, 126, 157, 200, 31, 115, 93, 182, 131, 195, 255, 110, 199, 177, 80, 63},
				parentFingerPrint: 3203769081,
				index:             2147483650,
				depth:             0x3,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromExtendedPublicKey(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromExtendedPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Error("object expected nil returned")
				t.FailNow()
			}

			if !assert.Equal(t, tt.want.PublicKey().Bytes(), got.PublicKey().Bytes()) {
				t.Errorf("fromExtendedPublicKey().PublicKey().Bytes() = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.chainCode, got.chainCode) {
				t.Errorf("fromExtendedPublicKey().chainCode = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.depth, got.depth) {
				t.Errorf("fromExtendedPublicKey().depth = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.index, got.index) {
				t.Errorf("fromExtendedPublicKey().index = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.parentFingerPrint, got.parentFingerPrint) {
				t.Errorf("fromExtendedPublicKey().parentFingerPrint = %v, want %v", got, tt.want)
			}
		})
	}
}
