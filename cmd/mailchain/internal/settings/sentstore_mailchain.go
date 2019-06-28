package settings

import (
	"github.com/mailchain/mailchain/stores"
)

type SentStoreMailchain struct {
	// domain string
}

func (m SentStoreMailchain) Produce() (*stores.SentStore, error) {
	return stores.NewSentStore(), nil
}
