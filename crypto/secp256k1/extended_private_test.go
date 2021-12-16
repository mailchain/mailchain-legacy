package secp256k1

import (
	"testing"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestExtendedPrivateKey_Bytes(t *testing.T) {
	type fields struct {
		key               PrivateKey
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
			"xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi",
			fields{
				func() PrivateKey {
					r, _ := PrivateKeyFromBytes(encodingtest.MustDecodeHex("e8f32e723decf4051aefac8e2c93c9c5b214313817cdb01a1494b917c8436b35"))
					return *r
				}(),
				[32]byte{135, 61, 255, 129, 192, 47, 82, 86, 35, 253, 31, 229, 22, 126, 172, 58, 85, 160, 73, 222, 61, 49, 75, 180, 46, 226, 39, 255, 237, 55, 213, 8},
				0,
				0,
				0,
			},
			[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x87, 0x3d, 0xff, 0x81, 0xc0, 0x2f, 0x52, 0x56, 0x23, 0xfd, 0x1f, 0xe5, 0x16, 0x7e, 0xac, 0x3a, 0x55, 0xa0, 0x49, 0xde, 0x3d, 0x31, 0x4b, 0xb4, 0x2e, 0xe2, 0x27, 0xff, 0xed, 0x37, 0xd5, 0x8, 0x0, 0xe8, 0xf3, 0x2e, 0x72, 0x3d, 0xec, 0xf4, 0x5, 0x1a, 0xef, 0xac, 0x8e, 0x2c, 0x93, 0xc9, 0xc5, 0xb2, 0x14, 0x31, 0x38, 0x17, 0xcd, 0xb0, 0x1a, 0x14, 0x94, 0xb9, 0x17, 0xc8, 0x43, 0x6b, 0x35},
		},
		{
			"depth-set",
			fields{
				func() PrivateKey {
					r, _ := PrivateKeyFromBytes(encodingtest.MustDecodeHex("e8f32e723decf4051aefac8e2c93c9c5b214313817cdb01a1494b917c8436b35"))
					return *r
				}(),
				[32]byte{135, 61, 255, 129, 192, 47, 82, 86, 35, 253, 31, 229, 22, 126, 172, 58, 85, 160, 73, 222, 61, 49, 75, 180, 46, 226, 39, 255, 237, 55, 213, 8},
				0,
				0,
				5,
			},
			[]byte{0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x87, 0x3d, 0xff, 0x81, 0xc0, 0x2f, 0x52, 0x56, 0x23, 0xfd, 0x1f, 0xe5, 0x16, 0x7e, 0xac, 0x3a, 0x55, 0xa0, 0x49, 0xde, 0x3d, 0x31, 0x4b, 0xb4, 0x2e, 0xe2, 0x27, 0xff, 0xed, 0x37, 0xd5, 0x8, 0x0, 0xe8, 0xf3, 0x2e, 0x72, 0x3d, 0xec, 0xf4, 0x5, 0x1a, 0xef, 0xac, 0x8e, 0x2c, 0x93, 0xc9, 0xc5, 0xb2, 0x14, 0x31, 0x38, 0x17, 0xcd, 0xb0, 0x1a, 0x14, 0x94, 0xb9, 0x17, 0xc8, 0x43, 0x6b, 0x35},
		},
		{
			"depth-fingerprint",
			fields{
				func() PrivateKey {
					r, _ := PrivateKeyFromBytes(encodingtest.MustDecodeHex("e8f32e723decf4051aefac8e2c93c9c5b214313817cdb01a1494b917c8436b35"))
					return *r
				}(),
				[32]byte{135, 61, 255, 129, 192, 47, 82, 86, 35, 253, 31, 229, 22, 126, 172, 58, 85, 160, 73, 222, 61, 49, 75, 180, 46, 226, 39, 255, 237, 55, 213, 8},
				4294961290,
				0,
				5,
			},
			[]byte{0x5, 0xff, 0xff, 0xe8, 0x8a, 0x0, 0x0, 0x0, 0x0, 0x87, 0x3d, 0xff, 0x81, 0xc0, 0x2f, 0x52, 0x56, 0x23, 0xfd, 0x1f, 0xe5, 0x16, 0x7e, 0xac, 0x3a, 0x55, 0xa0, 0x49, 0xde, 0x3d, 0x31, 0x4b, 0xb4, 0x2e, 0xe2, 0x27, 0xff, 0xed, 0x37, 0xd5, 0x8, 0x0, 0xe8, 0xf3, 0x2e, 0x72, 0x3d, 0xec, 0xf4, 0x5, 0x1a, 0xef, 0xac, 0x8e, 0x2c, 0x93, 0xc9, 0xc5, 0xb2, 0x14, 0x31, 0x38, 0x17, 0xcd, 0xb0, 0x1a, 0x14, 0x94, 0xb9, 0x17, 0xc8, 0x43, 0x6b, 0x35},
		},
		{
			"depth-fingerprint-index",
			fields{
				func() PrivateKey {
					r, _ := PrivateKeyFromBytes(encodingtest.MustDecodeHex("e8f32e723decf4051aefac8e2c93c9c5b214313817cdb01a1494b917c8436b35"))
					return *r
				}(),
				[32]byte{135, 61, 255, 129, 192, 47, 82, 86, 35, 253, 31, 229, 22, 126, 172, 58, 85, 160, 73, 222, 61, 49, 75, 180, 46, 226, 39, 255, 237, 55, 213, 8},
				4294961290,
				4294962854,
				5,
			},
			[]byte{0x5, 0xff, 0xff, 0xe8, 0x8a, 0xff, 0xff, 0xee, 0xa6, 0x87, 0x3d, 0xff, 0x81, 0xc0, 0x2f, 0x52, 0x56, 0x23, 0xfd, 0x1f, 0xe5, 0x16, 0x7e, 0xac, 0x3a, 0x55, 0xa0, 0x49, 0xde, 0x3d, 0x31, 0x4b, 0xb4, 0x2e, 0xe2, 0x27, 0xff, 0xed, 0x37, 0xd5, 0x8, 0x0, 0xe8, 0xf3, 0x2e, 0x72, 0x3d, 0xec, 0xf4, 0x5, 0x1a, 0xef, 0xac, 0x8e, 0x2c, 0x93, 0xc9, 0xc5, 0xb2, 0x14, 0x31, 0x38, 0x17, 0xcd, 0xb0, 0x1a, 0x14, 0x94, 0xb9, 0x17, 0xc8, 0x43, 0x6b, 0x35},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &ExtendedPrivateKey{
				key:               tt.fields.key,
				chainCode:         tt.fields.chainCode,
				parentFingerPrint: tt.fields.parentFingerPrint,
				index:             tt.fields.index,
				depth:             tt.fields.depth,
			}
			got := k.Bytes()

			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ExtendedPrivateKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromExtendedPrivateKey(t *testing.T) {
	type args struct {
		in *hdkeychain.ExtendedKey
	}
	tests := []struct {
		name    string
		args    args
		want    *ExtendedPrivateKey
		wantErr bool
	}{
		{
			"xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi",
			args{
				func() *hdkeychain.ExtendedKey {
					o, _ := hdkeychain.NewKeyFromString("xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi")
					return o
				}(),
			},
			&ExtendedPrivateKey{
				key: func() PrivateKey {
					o, _ := PrivateKeyFromBytes(encodingtest.MustDecodeHex("e8f32e723decf4051aefac8e2c93c9c5b214313817cdb01a1494b917c8436b35"))
					return *o
				}(),
				chainCode:         [32]byte{135, 61, 255, 129, 192, 47, 82, 86, 35, 253, 31, 229, 22, 126, 172, 58, 85, 160, 73, 222, 61, 49, 75, 180, 46, 226, 39, 255, 237, 55, 213, 8},
				parentFingerPrint: 0x0,
				index:             0x0,
				depth:             0x0,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromExtendedPrivateKey(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromExtendedPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Error("object expected nil returned")
				t.FailNow()
			}

			if !assert.Equal(t, tt.want.PrivateKey().Bytes(), got.PrivateKey().Bytes()) {
				t.Errorf("fromExtendedPrivateKey().PrivateKey().Bytes() = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.chainCode, got.chainCode) {
				t.Errorf("fromExtendedPrivateKey().chainCode = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.depth, got.depth) {
				t.Errorf("fromExtendedPrivateKey().depth = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.index, got.index) {
				t.Errorf("fromExtendedPrivateKey().index = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, tt.want.parentFingerPrint, got.parentFingerPrint) {
				t.Errorf("fromExtendedPrivateKey().parentFingerPrint = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestExtendedPrivateKey_ExtendedPublicKey(t *testing.T) {
	type fields struct {
		depth             byte
		parentFingerPrint uint32
		index             uint32
		chainCode         [32]byte
		key               PrivateKey
	}
	tests := []struct {
		name    string
		prvKey  crypto.ExtendedPrivateKey
		want    crypto.ExtendedPublicKey
		wantErr bool
	}{
		{
			"xprv9wNUHWVTuAHnj7y9JJRvdqgd8jsN5QuzdPt7EuBXfXXgjMEWPc5dENSs3HKvXvoPMyJsBpSMkEryBEz3kxdRg8fpAfq9RYh4wiysZihDR2r",
			func() crypto.ExtendedPrivateKey {
				hdKey, _ := hdkeychain.NewKeyFromString("xprv9wNUHWVTuAHnj7y9JJRvdqgd8jsN5QuzdPt7EuBXfXXgjMEWPc5dENSs3HKvXvoPMyJsBpSMkEryBEz3kxdRg8fpAfq9RYh4wiysZihDR2r")
				out, _ := fromExtendedPrivateKey(hdKey)
				return out
			}(),
			func() crypto.ExtendedPublicKey {
				hdKey, _ := hdkeychain.NewKeyFromString("xpub6AMph22MjXr5wc3cQKxvzydMgmhrUsdqzcoi3Hb9Ds4fc9Zew9PsnAmLtaBNTZCtzsZfLMgBM6DEFZGX2A4kHWDatJj6cfbRH896d2ACi4F")
				out, _ := fromExtendedPublicKey(hdKey)
				return out
			}(),
			false,
		},
		{
			"xprv9wHokC2KXdTSpEepFcu53hMDUHYfAtTaLEJEMyxBPAMf78hJg17WhL5FyeDUQH5KWmGjGgEb2j74gsZqgupWpPbZgP6uFmP8MYEy5BNbyET",
			func() crypto.ExtendedPrivateKey {
				hdKey, _ := hdkeychain.NewKeyFromString("xprv9wHokC2KXdTSpEepFcu53hMDUHYfAtTaLEJEMyxBPAMf78hJg17WhL5FyeDUQH5KWmGjGgEb2j74gsZqgupWpPbZgP6uFmP8MYEy5BNbyET")
				out, _ := fromExtendedPrivateKey(hdKey)
				return out
			}(),
			func() crypto.ExtendedPublicKey {
				hdKey, _ := hdkeychain.NewKeyFromString("xpub6AHA9hZDN11k2ijHMeS5QqHx2KP9aMBRhTDqANMnwVtdyw2TDYRmF8PjpvwUFcL1Et8Hj59S3gTSMcUQ5gAqTz3Wd8EsMTmF3DChhqPQBnU")
				out, _ := fromExtendedPublicKey(hdKey)
				return out
			}(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.prvKey.ExtendedPublicKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtendedPrivateKey.ExtendedPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ExtendedPrivateKey.ExtendedPublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
