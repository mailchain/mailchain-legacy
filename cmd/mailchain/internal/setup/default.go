package setup

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
	"github.com/spf13/viper" // nolint: depguard
)

func DefaultKeystore() Keystore {
	return Keystore{
		setter:             config.DefaultKeystore(),
		selectItemSkipable: prompts.SelectItemSkipable,
		viper:              viper.GetViper(),
	}
}

func DefaultNetwork() Network {
	return Network{
		receiverSelector:     DefaultReceiver(),
		senderSelector:       DefaultSender(),
		pubKeyFinderSelector: DefaultPubKeyFinder(),
		selectItem:           prompts.SelectItem,
	}
}

func DefaultPubKeyFinder() PubKeyFinder {
	return PubKeyFinder{
		setter:             config.DefaultPubKeyFinder(),
		selectItemSkipable: prompts.SelectItemSkipable,
		viper:              viper.GetViper(),
	}
}

func DefaultReceiver() Receiver {
	return Receiver{
		setter:             config.DefaultReceiver(),
		selectItemSkipable: prompts.SelectItemSkipable,
		viper:              viper.GetViper(),
	}
}

func DefaultSender() Sender {
	return Sender{
		setter:             config.DefaultSender(),
		selectItemSkipable: prompts.SelectItemSkipable,
		viper:              viper.GetViper(),
	}
}
func DefaultSentStorage() SentStorage {
	return SentStorage{
		setter:             config.DefaultSentStore(),
		selectItemSkipable: prompts.SelectItemSkipable,
		viper:              viper.GetViper(),
	}
}
