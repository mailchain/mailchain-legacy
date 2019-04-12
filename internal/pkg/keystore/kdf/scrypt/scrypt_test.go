package scrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDefault(t *testing.T) {
	assert := assert.New(t)
	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions()})
	assert.Equal(32, opts.Len)
	assert.Equal(262144, opts.N)
	assert.Equal(1, opts.P)
	assert.Equal("", opts.Passphrase)
	assert.Equal(8, opts.R)
	assert.Nil(opts.Salt)
}

func TestApplyDefaultAndPassword(t *testing.T) {
	assert := assert.New(t)
	randomSalt, err := RandomSalt()
	if err != nil {
		t.Fail()
	}

	opts := &DeriveOpts{}
	apply(opts, []DeriveOptionsBuilder{DefaultDeriveOptions(), WithPassphrase("test"), randomSalt})
	assert.Equal(32, opts.Len)
	assert.Equal(262144, opts.N)
	assert.Equal(1, opts.P)
	assert.Equal("test", opts.Passphrase)
	assert.Equal(8, opts.R)
	assert.Equal(32, len(opts.Salt))
}
