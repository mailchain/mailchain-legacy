package config

import (
	"github.com/imdario/mergo"
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

func DefaultClients() *Clients {
	return &Clients{
		viper:         viper.GetViper(),
		requiredInput: prompts.RequiredInput,
	}
}

func DefaultPubKeyFinder() *PubKeyFinder {
	clients := DefaultClients()
	return &PubKeyFinder{
		viper:        viper.GetViper(),
		clientGetter: clients,
		clientSetter: clients,
		mapMerge:     mergo.Merge,
	}
}
