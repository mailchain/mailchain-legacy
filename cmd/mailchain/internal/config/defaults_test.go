package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSentStore(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want SentStore
	}{
		{
			"success",
			SentStore{
				viper: viper.GetViper(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultSentStore()
			if !assert.EqualValues(tt.want.viper, got.viper) {
				t.Errorf("DefaultSentStore().viper = %v, want %v", got.viper, tt.want.viper)
			}
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.requiredInput) {
				t.Error("want got.requiredInput != nil")
			}
		})
	}
}

func TestDefaultKeystore(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want *Keystore
	}{
		{
			"success",
			&Keystore{
				viper: viper.GetViper(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultKeystore()
			if !assert.EqualValues(tt.want.viper, got.viper) {
				t.Errorf("DefaultKeystore().viper = %v, want %v", got.viper, tt.want.viper)
			}
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.requiredInputWithDefault) {
				t.Error("want got.requiredInputWithDefault != nil")
			}
		})
	}
}

func TestDefaultClients(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want *Clients
	}{
		{
			"success",
			&Clients{
				viper: viper.GetViper(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultClients()
			if !assert.EqualValues(tt.want.viper, got.viper) {
				t.Errorf("DefaultClients().viper = %v, want %v", got.viper, tt.want.viper)
			}
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}
			if !assert.NotNil(got.requiredInput) {
				t.Error("want got.requiredInput != nil")
			}
		})
	}
}

func TestDefaultPubKeyFinder(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want *PubKeyFinder
	}{
		{
			"success",
			&PubKeyFinder{
				viper: viper.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultPubKeyFinder()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}

			if !assert.EqualValues(tt.want.viper, got.viper) {
				t.Errorf("DefaultPubKeyFinder().viper = %v, want %v", got.viper, tt.want.viper)
			}
			if !assert.NotNil(got.mapMerge) {
				t.Error("want got.mapMerge != nil")
			}
			if !assert.IsType(&Clients{}, got.clientGetter) {
				t.Error("invalid clientGetter type")
			}
			if !assert.IsType(&Clients{}, got.clientSetter) {
				t.Error("invalid clientGetter type")
			}
		})
	}
}

func TestDefaultReceiver(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want *Receiver
	}{
		{
			"success",
			&Receiver{
				viper: viper.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultReceiver()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}

			if !assert.EqualValues(tt.want.viper, got.viper) {
				t.Errorf("DefaultReceiver().viper = %v, want %v", got.viper, tt.want.viper)
			}
			if !assert.NotNil(got.mapMerge) {
				t.Error("want got.mapMerge != nil")
			}
			if !assert.IsType(&Clients{}, got.clientGetter) {
				t.Error("invalid clientGetter type")
			}
			if !assert.IsType(&Clients{}, got.clientSetter) {
				t.Error("invalid clientGetter type")
			}
		})
	}
}

func TestDefaultSender(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want *Sender
	}{
		{
			"success",
			&Sender{
				viper: viper.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultSender()
			if !assert.NotNil(got) {
				t.Error("want got != nil")
			}

			if !assert.EqualValues(tt.want.viper, got.viper) {
				t.Errorf("DefaultSender().viper = %v, want %v", got.viper, tt.want.viper)
			}
			if !assert.NotNil(got.mapMerge) {
				t.Error("want got.mapMerge != nil")
			}
			if !assert.IsType(&Clients{}, got.clientGetter) {
				t.Error("invalid clientGetter type")
			}
			if !assert.IsType(&Clients{}, got.clientSetter) {
				t.Error("invalid clientGetter type")
			}
		})
	}
}
