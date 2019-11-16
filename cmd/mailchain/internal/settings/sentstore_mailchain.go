package settings

import (
	"github.com/mailchain/mailchain/stores"
)

// SentStoreMailchain configuration element.
type SentStoreMailchain struct {
	// domain string
}

// Produce `stores.SentStore` based on configuration settings.
func (m SentStoreMailchain) Produce() (*stores.SentStore, error) {
	return stores.NewSentStore(), nil
}
