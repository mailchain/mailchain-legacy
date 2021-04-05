package integration

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/algorand"
	"github.com/mailchain/mailchain/stores"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSendReceive(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	v := viper.New()
	v.SetConfigFile("./private-keys.yaml")

	if err := v.ReadInConfig(); !assert.NoError(t, err) {
		t.FailNow()
	}

	type fields struct {
		settings map[string]interface{}
	}

	type args struct {
		protocol      string
		network       string
		fromKeyLookup string
		toKeyLookup   string
		contentType   string
		envelope      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []stores.Transaction
		wantErr bool
	}{
		{
			"success-algorand-testnet-send-charlotte-sofia",
			fields{},
			args{
				protocols.Algorand,
				algorand.Testnet,
				"algorand.testnet.charlotte",
				"algorand.testnet.sofia",
				"'text/plain; charset=\\\"UTF-8\\\"'",
				"0x01",
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passphrase := "int-test-secret"
			testDir := testDir(t)

			if err := os.MkdirAll(testDir, os.ModePerm); !assert.NoError(t, err) {
				t.FailNow()
			}
			if err := os.MkdirAll(fmt.Sprintf("%s/api", testDir), os.ModePerm); !assert.NoError(t, err) {
				t.FailNow()
			}

			settingsFileName, err := createTestCaseSettingsFile(tt.fields.settings, testDir)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			toPubkey := addPrivateKey(t, tt.args.protocol, tt.args.network, settingsFileName, passphrase, v, tt.args.toKeyLookup)
			fromPubKey := addPrivateKey(t, tt.args.protocol, tt.args.network, settingsFileName, passphrase, v, tt.args.fromKeyLookup)

			serveCmd := createCommand(settingsFileName, bundle("serve"), bundle("--passphrase", passphrase))

			var serveOutBuffer bytes.Buffer
			serveCmd.Stdout = &serveOutBuffer
			serveCmd.Stderr = &serveOutBuffer

			if err := serveCmd.Start(); err != nil {
				t.FailNow()
			}
			t.Cleanup(func() {
				t.Logf("output: serve\n%s", serveOutBuffer.Bytes())
				syscall.Kill(-serveCmd.Process.Pid, syscall.SIGKILL)
			})

			time.Sleep(10 * time.Second)

			toAddress := encodeAddress(t, toPubkey, tt.args.protocol, tt.args.network)
			fromAddress := encodeAddress(t, fromPubKey, tt.args.protocol, tt.args.network)

			apiCheckContainsAddress(t, fromAddress, tt.args.protocol, tt.args.network)
			apiCheckContainsAddress(t, toAddress, tt.args.protocol, tt.args.network)

			toPubKeyRes := apiGetPublicKey(t, toAddress, tt.args.protocol, tt.args.network)

			subject := apiSendMessage(t, tt.args.protocol, tt.args.network, tt.args.contentType, tt.args.envelope, toPubKeyRes.SupportedEncryptionTypes[0], toAddress, fromAddress, toPubkey)

			time.Sleep(30 * time.Second)

			apiCheckMessage(t, tt.args.protocol, tt.args.network, toAddress, subject)
		})
	}
}
