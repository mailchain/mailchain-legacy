package integration

import (
	"os"
	"strings"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/internal/addressing"
	"github.com/stretchr/testify/assert"
)

func testDir(t *testing.T) string {
	wd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return strings.Join([]string{wd, "out", t.Name()}, "/")
}

func encodeAddress(t *testing.T, pubKey crypto.PublicKey, protocol, network string) string {
	addressBytes, err := addressing.FromPublicKey(pubKey, protocol, network)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	encoded, _, err := addressing.EncodeByProtocol(addressBytes, protocol)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	return encoded
}
