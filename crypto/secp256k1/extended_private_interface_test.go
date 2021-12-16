package secp256k1_test

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestExtendedPrivateKey_RoundTripSerialization(t *testing.T) {
	tests := []struct {
		name      string
		coreBytes []byte
		wantErr   bool
	}{
		{
			"1",
			encodingtest.MustDecodeBase58("xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi")[4:78],
			false,
		},
		{
			"2",
			encodingtest.MustDecodeBase58("xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U")[4:78],
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := secp256k1.ExtendedPrivateKeyFromBytes(tt.coreBytes) // only need the relevant bytes
			if !assert.NoError(t, err) {
				assert.FailNow(t, "failed to unmarshal key")
			}

			if !assert.Equal(t, k.Bytes(), tt.coreBytes) {
				t.Errorf("Roundtrip failed coreBytes = %v, k.Bytes() %v", tt.coreBytes, k.Bytes())
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func privateKeyFromBIP32String(t *testing.T, in string) crypto.ExtendedPrivateKey {
	b, err := encoding.DecodeBase58(in)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "failed to decode input")
	}

	k, err := secp256k1.ExtendedPrivateKeyFromBytes(b[4:78]) // only need the relevant bytes
	if !assert.NoError(t, err) {
		assert.FailNow(t, "failed to unmarshal key")
	}

	return k
}
func TestExtendedPrivateKey_Derive(t *testing.T) {
	// The private extended keys for test vectors in [BIP32].
	testVec1MasterPrivKey := "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"
	testVec2MasterPrivKey := "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U"

	tests := []struct {
		name    string
		baseKey crypto.ExtendedPrivateKey
		path    []uint32
		want    crypto.ExtendedPrivateKey
		wantErr bool
	}{
		{
			"test vector 1 chain m",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			[]uint32{},
			privateKeyFromBIP32String(t, "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"),
			false,
		},
		{
			"test vector 1 chain m/0",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			[]uint32{0},
			privateKeyFromBIP32String(t, "xprv9uHRZZhbkedL37eZEnyrNsQPFZYRAvjy5rt6M1nbEkLSo378x1CQQLo2xxBvREwiK6kqf7GRNvsNEchwibzXaV6i5GcsgyjBeRguXhKsi4R"),
			false,
		},
		{
			"test vector 1 chain m/0/1",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			[]uint32{0, 1},
			privateKeyFromBIP32String(t, "xprv9ww7sMFLzJMzy7bV1qs7nGBxgKYrgcm3HcJvGb4yvNhT9vxXC7eX7WVULzCfxucFEn2TsVvJw25hH9d4mchywguGQCZvRgsiRaTY1HCqN8G"),
			false,
		},
		{
			"test vector 1 chain m/0/1/2",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			[]uint32{0, 1, 2},
			privateKeyFromBIP32String(t, "xprv9xrdP7iD2L1YZCgR9AecDgpDMZSTzP5KCfUykGXgjBxLgp1VFHsEeL3conzGAkbc1MigG1o8YqmfEA2jtkPdf4vwMaGJC2YSDbBTPAjfRUi"),
			false,
		},
		{
			"test vector 1 chain m/0/1/2/2",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			[]uint32{0, 1, 2, 2},
			privateKeyFromBIP32String(t, "xprvA2J8Hq4eiP7xCEBP7gzRJGJnd9CHTkEU6eTNMrZ6YR7H5boik8daFtDZxmJDfdMSKHwroCfAfsBKWWidRfBQjpegy6kzXSkQGGoMdWKz5Xh"),
			false,
		},
		{
			"test vector 1 chain m/0/1/2/2/1000000000",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			[]uint32{0, 1, 2, 2, 1000000000},
			privateKeyFromBIP32String(t, "xprvA3XhazxncJqJsQcG85Gg61qwPQKiobAnWjuPpjKhExprZjfse6nErRwTMwGe6uGWXPSykZSTiYb2TXAm7Qhwj8KgRd2XaD21Styu6h6AwFz"),
			false,
		},

		// // Test vector 2
		{
			"test vector 2 chain m",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			[]uint32{},
			privateKeyFromBIP32String(t, "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U"),
			false,
		},
		{
			"test vector 2 chain m/0",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			[]uint32{0},
			privateKeyFromBIP32String(t, "xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt"),
			false,
		},
		{
			"test vector 2 chain m/0/2147483647",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			[]uint32{0, 2147483647},
			privateKeyFromBIP32String(t, "xprv9wSp6B7cXJWXZRpDbxkFg3ry2fuSyUfvboJ5Yi6YNw7i1bXmq9QwQ7EwMpeG4cK2pnMqEx1cLYD7cSGSCtruGSXC6ZSVDHugMsZgbuY62m6"),
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			[]uint32{0, 2147483647, 1},
			privateKeyFromBIP32String(t, "xprv9ysS5br6UbWCRCJcggvpUNMyhVWgD7NypY9gsVTMYmuRtZg8izyYC5Ey4T931WgWbfJwRDwfVFqV3b29gqHDbuEpGcbzf16pdomk54NXkSm"),
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1/2147483646",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			[]uint32{0, 2147483647, 1, 2147483646},
			privateKeyFromBIP32String(t, "xprvA2LfeWWwRCxh4iqigcDMnUf2E3nVUFkntc93nmUYBtb9rpSPYWa8MY3x9ZHSLZkg4G84UefrDruVK3FhMLSJsGtBx883iddHNuH1LNpRrEp"),
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1/2147483646/2",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			[]uint32{0, 2147483647, 1, 2147483646, 2},
			privateKeyFromBIP32String(t, "xprvA48ALo8BDjcRET68R5RsPzF3H7WeyYYtHcyUeLRGBPHXu6CJSGjwW7dWoeUWTEzT7LG3qk6Eg6x2ZoqD8gtyEFZecpAyvchksfLyg3Zbqam"),
			false,
		},

		// Custom tests to trigger specific conditions.
		{
			// Seed 000000000000000000000000000000da.
			"Derived privkey with zero high byte m/0",
			privateKeyFromBIP32String(t, "xprv9s21ZrQH143K4FR6rNeqEK4EBhRgLjWLWhA3pw8iqgAKk82ypz58PXbrzU19opYcxw8JDJQF4id55PwTsN1Zv8Xt6SKvbr2KNU5y8jN8djz"),
			[]uint32{0},
			privateKeyFromBIP32String(t, "xprv9uC5JqtViMmgcAMUxcsBCBFA7oYCNs4bozPbyvLfddjHou4rMiGEHipz94xNaPb1e4f18TRoPXfiXx4C3cDAcADqxCSRSSWLvMBRWPctSN9"),
			false,
		},

		{
			"err-too-deep",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			[]uint32{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extKey := tt.baseKey
			var err error

			for _, childNum := range tt.path {
				extKey, err = extKey.Derive(childNum)
				if err != nil {
					break
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Derive error = %v, wantErr %v", err, tt.wantErr)
				t.FailNow()
			}

			if extKey == nil || tt.want == nil {
				return
			}

			if !assert.Equal(t, tt.want.Bytes(), extKey.Bytes()) {
				t.Errorf("ExtendedPrivateKey.Derive() = %v, want %v", extKey, tt.want)
			}
		})
	}
}
