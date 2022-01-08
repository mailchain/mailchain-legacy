package bip32_test

import (
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestPrivatePublic_Derive(t *testing.T) {
	testVec1MasterPrivKey := "xprv9wNUHWVTuAHnj7y9JJRvdqgd8jsN5QuzdPt7EuBXfXXgjMEWPc5dENSs3HKvXvoPMyJsBpSMkEryBEz3kxdRg8fpAfq9RYh4wiysZihDR2r"
	testVec1MasterPubKey := "xpub6AMph22MjXr5wc3cQKxvzydMgmhrUsdqzcoi3Hb9Ds4fc9Zew9PsnAmLtaBNTZCtzsZfLMgBM6DEFZGX2A4kHWDatJj6cfbRH896d2ACi4F"
	testVec2MasterPrivKey := "xprv9wHokC2KXdTSpEepFcu53hMDUHYfAtTaLEJEMyxBPAMf78hJg17WhL5FyeDUQH5KWmGjGgEb2j74gsZqgupWpPbZgP6uFmP8MYEy5BNbyET"
	testVec2MasterPubKey := "xpub6AHA9hZDN11k2ijHMeS5QqHx2KP9aMBRhTDqANMnwVtdyw2TDYRmF8PjpvwUFcL1Et8Hj59S3gTSMcUQ5gAqTz3Wd8EsMTmF3DChhqPQBnU"

	tests := []struct {
		name       string
		privateKey crypto.ExtendedPrivateKey
		publicKey  crypto.ExtendedPublicKey
		path       []uint32
		wantErr    bool
	}{
		{
			"test vector 1 chain m/0/1",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1},
			false,
		},
		{
			"test vector 1 chain m",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{},
			false,
		},
		{
			"test vector 1 chain m/0",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0},
			false,
		},
		{
			"test vector 1 chain m/0/1",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1},
			false,
		},
		{
			"test vector 1 chain m/0/1/2",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1, 2},
			false,
		},
		{
			"test vector 1 chain m/0/1/2/2",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1, 2, 2},
			false,
		},
		{
			"test vector 1 chain m/0/1/2/2/1000000000",
			privateKeyFromBIP32String(t, testVec1MasterPrivKey),
			publicKeyFromBIP32String(t, testVec1MasterPubKey),
			[]uint32{0, 1, 2, 2, 1000000000},
			false,
		},

		// // Test vector 2
		{
			"test vector 2 chain m",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{},
			false,
		},
		{
			"test vector 2 chain m/0",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0},
			false,
		},
		{
			"test vector 2 chain m/0/2147483647",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0, 2147483647},
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0, 2147483647, 1},
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1/2147483646",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0, 2147483647, 1, 2147483646},
			false,
		},
		{
			"test vector 2 chain m/0/2147483647/1/2147483646/2",
			privateKeyFromBIP32String(t, testVec2MasterPrivKey),
			publicKeyFromBIP32String(t, testVec2MasterPubKey),
			[]uint32{0, 2147483647, 1, 2147483646, 2},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extPrvKey := tt.privateKey
			extPubKey := tt.publicKey
			var err error

			for _, childNum := range tt.path {
				extPrvKey, err = extPrvKey.Derive(childNum)
				if !assert.NoError(t, err) {
					assert.FailNow(t, "error not expected")
				}
			}

			for _, childNum := range tt.path {
				extPubKey, err = extPubKey.Derive(childNum)
				if !assert.NoError(t, err) {
					assert.FailNow(t, "error not expected")
				}
			}

			pubKey, err := extPrvKey.ExtendedPublicKey()
			if !assert.NoError(t, err) {
				assert.FailNow(t, "error not expected")
			}

			if !assert.Equal(t, pubKey, extPubKey) {
				t.Errorf("RoundTripParse and Derive = %v, want %v", pubKey.Bytes(), extPubKey.Bytes())
			}
		})
	}
}
