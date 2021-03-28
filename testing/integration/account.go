package integration

import (
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/mailchain/mailchain/encoding"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func addPrivateKey(t *testing.T, protocol, network, settingsFileName, passphrase string, v *viper.Viper, keyLookup string) crypto.PublicKey {
	privateKey := v.GetString(keyLookup + ".private-key")
	if privateKey == "" {
		assert.FailNow(t, "private key not set")
	}
	privateKeyEncoding := v.GetString(keyLookup + ".private-key-encoding")
	if privateKey == "" {
		assert.FailNow(t, "private key encoding not set")
	}
	keyType := v.GetString(keyLookup + ".key-type")
	if privateKey == "" {
		assert.FailNow(t, "key type not set")
	}

	createAccountCmd := func(privateKey, privateKeyEncoding, keyType string) *exec.Cmd {
		return createCommand(settingsFileName, bundle("account", "add"),
			bundle(
				"--private-key", privateKey,
				"--private-key-encoding", privateKeyEncoding,
				"--key-type", keyType,
				"--passphrase", passphrase,
				"--protocol", protocol,
				"--network", network,
			),
		)
	}

	out, err := createAccountCmd(privateKey, privateKeyEncoding, keyType).CombinedOutput()
	if !assert.NoError(t, err) {
		t.Logf("failed to add account: %s", out)
		return nil
	}

	type response struct {
		Message           string `json:"message"`
		PublicKey         string `json:"public-key"`
		PublicKeyEncoding string `json:"public-key-encoding"`
		Protocol          string `json:"protocol"`
		Network           string `json:"network"`
	}

	var cmdRes response

	if err := json.Unmarshal(out, &cmdRes); !assert.NoError(t, err) {
		t.FailNow()
	}

	pubKeyBytes, err := encoding.Decode(cmdRes.PublicKeyEncoding, cmdRes.PublicKey)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	pubKey, err := multikey.PublicKeyFromBytes(keyType, pubKeyBytes)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	return pubKey
}
