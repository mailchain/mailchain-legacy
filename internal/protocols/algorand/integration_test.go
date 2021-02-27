package algorand_test

import (
	"testing"

	"github.com/mailchain/mailchain/crypto/ed25519"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/stretchr/testify/assert"
)

func TestAddPrivateKeyFromMnemonic(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name          string
		args          args
		want          []byte
		wantPublicKey string
		wantAddress   string
		wantErr       bool
	}{
		{
			"success-25-words",
			args{
				"subject woman doll exercise order intact sting crawl course shallow provide keen lounge dog velvet immune ethics hour emotion defense guitar second local absent bullet",
			},
			[]byte{0xbf, 0x3e, 0x7f, 0x81, 0xf4, 0x14, 0x4e, 0xd6, 0xbd, 0xba, 0x32, 0x8a, 0x41, 0xf1, 0x59, 0x9b, 0x27, 0x42, 0x2, 0x41, 0xbe, 0x71, 0x6d, 0x92, 0x9b, 0x91, 0x96, 0xe3, 0xb3, 0x9, 0x6b, 0xb0},
			"VCPSID7JBX2R252YHDVJMMKELH6WB3WYGPGQEZ4MYU3T7XKI2VNQ",
			"VCPSID7JBX2R252YHDVJMMKELH6WB3WYGPGQEZ4MYU3T7XKI2VNUXKOP7A",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encoding.DecodeMnemonicAlgorand(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeMnemonicAlgorand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("DecodeMnemonicAlgorand() = %v, want %v", got, tt.want)
			}
			pk, _ := ed25519.PrivateKeyFromBytes(got)
			gotPubKey := pk.PublicKey()
			if !assert.Equal(t, tt.wantPublicKey, encoding.EncodeBase32(gotPubKey.Bytes())) {
				t.Errorf("PublicKey = %v, want %v", encoding.EncodeBase32(gotPubKey.Bytes()), tt.wantPublicKey)
			}

			addBytes, err := algorand.Address(gotPubKey)

			gotAddress := encoding.EncodeBase32(addBytes)
			if !assert.Equal(t, tt.wantAddress, gotAddress) {
				t.Errorf("Address = %v, want %v", gotAddress, tt.wantAddress)
			}
			assert.NoError(t, err)
		})
	}
}
