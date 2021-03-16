package algod

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New algod client is created that stores config for each supported network.
func New(algodAddress, algodToken string) (*Client, error) {
	return &Client{
		algoAddress: algodAddress,
		algodToken:  algodToken,
		logger:      log.With().Str("component", "sender").Str("kind", "algod").Logger(),
	}, nil
}

// Client ethereum JSON-RPC2 client that is used to send transactions.
type Client struct {
	algoAddress string
	algodToken  string
	logger      zerolog.Logger
}

type networkConfig struct {
	url string
}
