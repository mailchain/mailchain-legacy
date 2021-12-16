package secp256k1_test

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/stretchr/testify/assert"
)

func TestExtendedPublicKey_RoundTripSerialization(t *testing.T) {
	tests := []struct {
		name      string
		coreBytes []byte
		wantErr   bool
	}{
		{
			"1",
			encodingtest.MustDecodeBase58("xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8")[4:78],
			false,
		},
		{
			"2",
			encodingtest.MustDecodeBase58("xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB")[4:78],
			false,
		},
		{
			"3",
			encodingtest.MustDecodeBase58("xpub6AvUGrnEpfvJBbfx7sQ89Q8hEMPM65UteqEX4yUbUiES2jHfjexmfJoxCGSwFMZiPBaKQT1RiKWrKfuDV4vpgVs4Xn8PpPTR2i79rwHd4Zr")[4:78],
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := secp256k1.ExtendedPublicKeyFromBytes(tt.coreBytes) // only need the relevant bytes
			if !assert.NoError(t, err) {
				assert.FailNow(t, "failed to unmarshal key")
			}

			if !assert.Equal(t, tt.coreBytes, k.Bytes()) {
				t.Errorf("Roundtrip failed coreBytes = %v, k.Bytes() %v", tt.coreBytes, k.Bytes())
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func publicKeyFromBIP32String(t *testing.T, in string) crypto.ExtendedPublicKey {
	b, err := encoding.DecodeBase58(in)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "failed to decode input")
	}

	k, err := secp256k1.ExtendedPublicKeyFromBytes(b[4:78]) // only need the relevant bytes
	if !assert.NoError(t, err) {
		assert.FailNow(t, "failed to unmarshal key")
	}

	return k
}

func TestExtendedPublicKey_Derive(t *testing.T) {
	// The public extended keys for test vectors in [BIP32].
	testVec1MasterPubKey := "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"
	testVec2MasterPubKey := "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB"

	tests := []struct {
		name    string
		baseKey crypto.ExtendedPublicKey
		path    []uint32
		want    crypto.ExtendedPublicKey
		wantErr bool
	}{
		{
			"test vector 1 chain m",
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{},
			publicKeyFromBIP32String(t, "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"),
			false,
		},
		{
			"test vector 1 chain m/0",
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0},
			publicKeyFromBIP32String(t, "xpub68Gmy5EVb2BdFbj2LpWrk1M7obNuaPTpT5oh9QCCo5sRfqSHVYWex97WpDZzszdzHzxXDAzPLVSwybe4uPYkSk4G3gnrPqqkV9RyNzAcNJ1"),
			false,
		},
		{
			"test vector 1 chain m/0/1",
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1},
			publicKeyFromBIP32String(t, "xpub6AvUGrnEpfvJBbfx7sQ89Q8hEMPM65UteqEX4yUbUiES2jHfjexmfJoxCGSwFMZiPBaKQT1RiKWrKfuDV4vpgVs4Xn8PpPTR2i79rwHd4Zr"),
			false,
		},
		{
			"test vector 1 chain m/0/1/2",
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1, 2},
			publicKeyFromBIP32String(t, "xpub6BqyndF6rhZqmgktFCBcapkwubGxPqoAZtQaYewJHXVKZcLdnqBVC8N6f6FSHWUghjuTLeubWyQWfJdk2G3tGgvgj3qngo4vLTnnSjAZckv"),
			false,
		},
		{
			"test vector 1 chain m/0/1/2/2",
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1, 2, 2},
			publicKeyFromBIP32String(t, "xpub6FHUhLbYYkgFQiFrDiXRfQFXBB2msCxKTsNyAExi6keFxQ8sHfwpogY3p3s1ePSpUqLNYks5T6a3JqpCGszt4kxbyq7tUoFP5c8KWyiDtPp"),
			false,
		},
		{
			"test vector 1 chain m/0/1/2/2/1000000000",
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1, 2, 2, 1000000000},
			publicKeyFromBIP32String(t, "xpub6GX3zWVgSgPc5tgjE6ogT9nfwSADD3tdsxpzd7jJoJMqSY12Be6VQEFwDCp6wAQoZsH2iq5nNocHEaVDxBcobPrkZCjYW3QUmoDYzMFBDu9"),
			false,
		},

		// Test vector 2
		{
			"test vector 2 chain m",
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{},
			publicKeyFromBIP32String(t, "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB"),
			false,
		},
		{
			"test vector 2 chain m/0",
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0},
			publicKeyFromBIP32String(t, "xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH"),
			false,
		},
		{
			"test vector 2 chain m/0/2147483647",
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0, 2147483647},
			publicKeyFromBIP32String(t, "xpub6ASAVgeWMg4pmutghzHG3BohahjwNwPmy2DgM6W9wGegtPrvNgjBwuZRD7hSDFhYfunq8vDgwG4ah1gVzZysgp3UsKz7VNjCnSUJJ5T4fdD"),
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1",
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0, 2147483647, 1},
			publicKeyFromBIP32String(t, "xpub6CrnV7NzJy4VdgP5niTpqWJiFXMAca6qBm5Hfsry77SQmN1HGYHnjsZSujoHzdxf7ZNK5UVrmDXFPiEW2ecwHGWMFGUxPC9ARipss9rXd4b"),
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1/2147483646",
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0, 2147483647, 1, 2147483646},
			publicKeyFromBIP32String(t, "xpub6FL2423qFaWzHCvBndkN9cbkn5cysiUeFq4eb9t9kE88jcmY63tNuLNRzpHPdAM4dUpLhZ7aUm2cJ5zF7KYonf4jAPfRqTMTRBNkQL3Tfta"),
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1/2147483646/2",
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0, 2147483647, 1, 2147483646, 2},
			publicKeyFromBIP32String(t, "xpub6H7WkJf547AiSwAbX6xsm8Bmq9M9P1Gjequ5SipsjipWmtXSyp4C3uwzewedGEgAMsDy4jEvNTWtxLyqqHY9C12gaBmgUdk2CGmwachwnWK"),
			false,
		},

		{
			"err-too-deep",
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
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
				t.Errorf("ExtendedPublicKey.Derive() = %v, want %v", extKey, tt.want)
			}
		})
	}
}
