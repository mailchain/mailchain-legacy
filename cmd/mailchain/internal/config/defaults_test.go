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
