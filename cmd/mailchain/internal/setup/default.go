package setup

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
	"github.com/spf13/viper"
)

func DefaultSentStorage() SentStorage {
	return SentStorage{
		sentStoreSetter:    config.DefaultSentStore(),
		selectItemSkipable: prompts.SelectItemSkipable,
		viper:              viper.GetViper(),
	}
}

func DefaultKeystore() Keystore {
	return Keystore{
		keystoreSetter:     config.DefaultKeystore(),
		selectItemSkipable: prompts.SelectItemSkipable,
		viper:              viper.GetViper(),
	}
}
