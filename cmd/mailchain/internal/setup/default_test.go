package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSentStorage(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
	}{
		{
			"success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultSentStorage()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.selectItemSkipable) {
				t.Error("want got.selectItemSkipable != nil")
			}
			if !assert.NotNil(got.setter) {
				t.Error("want got.setter != nil")
			}
			if !assert.NotNil(got.viper) {
				t.Error("want got.viper != nil")
			}
		})
	}
}

func TestDefaultKeystore(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
	}{
		{
			"success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultKeystore()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.selectItemSkipable) {
				t.Error("want got.selectItemSkipable != nil")
			}
			if !assert.NotNil(got.setter) {
				t.Error("want got.setter != nil")
			}
			if !assert.NotNil(got.viper) {
				t.Error("want got.viper != nil")
			}
		})
	}
}

func TestDefaultSender(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
	}{
		{"success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultSender()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.selectItemSkipable) {
				t.Error("want got.selectItemSkipable != nil")
			}
			if !assert.NotNil(got.setter) {
				t.Error("want got.setter != nil")
			}
			if !assert.NotNil(got.viper) {
				t.Error("want got.viper != nil")
			}
		})
	}
}

func TestDefaultPubKeyFinder(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
	}{
		{"success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultPubKeyFinder()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.selectItemSkipable) {
				t.Error("want got.selectItemSkipable != nil")
			}
			if !assert.NotNil(got.setter) {
				t.Error("want got.setter != nil")
			}
			if !assert.NotNil(got.viper) {
				t.Error("want got.viper != nil")
			}
		})
	}
}

func TestDefaultReceiver(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
	}{
		{"success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultReceiver()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.selectItemSkipable) {
				t.Error("want got.selectItemSkipable != nil")
			}
			if !assert.NotNil(got.setter) {
				t.Error("want got.setter != nil")
			}
			if !assert.NotNil(got.viper) {
				t.Error("want got.viper != nil")
			}
		})
	}
}

func TestDefaultNetwork(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
	}{
		{"success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultNetwork()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.selectItem) {
				t.Error("want got.selectItem != nil")
			}
			if !assert.NotNil(got.receiverSelector) {
				t.Error("want got.receiverSelector != nil")
			}
			if !assert.NotNil(got.senderSelector) {
				t.Error("want got.senderSelector != nil")
			}
			if !assert.NotNil(got.pubKeyFinderSelector) {
				t.Error("want got.pubKeyFinderSelector != nil")
			}
		})
	}
}
