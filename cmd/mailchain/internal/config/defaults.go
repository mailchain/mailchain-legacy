package config

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
	"github.com/spf13/viper" // nolint: depguard
)

func DefaultSentStore() *SentStore {
	return &SentStore{
		viper:         viper.GetViper(),
		requiredInput: prompts.RequiredInput,
	}
}

func DefaultKeystore() *Keystore {
	return &Keystore{
		viper:                    viper.GetViper(),
		requiredInputWithDefault: prompts.RequiredInputWithDefault,
	}
}
