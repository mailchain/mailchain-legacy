package setup

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/config"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/prompts"
)

func DefaultSentStorage() SentStorage {
	return SentStorage{
		sentStoreSetter:    config.DefaultSentStore(),
		selectItemSkipable: prompts.SelectItemSkipable,
	}
}
