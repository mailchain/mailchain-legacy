package commands

import (
	"context"
	"strconv"

	"github.com/mailchain/mailchain/cmd/indexer/internal/clients"
)

func getBlockNumber(blockNumber string, client clients.Block) (uint64, error) {
	if blockNumber == "latest" {
		return client.LatestBlockNumber(context.Background())
	}

	return strconv.ParseUint(blockNumber, 0, 64)
}
