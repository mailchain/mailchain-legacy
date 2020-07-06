package commands

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/indexer/internal/processor"
	sub "github.com/mailchain/mailchain/cmd/indexer/internal/substrate"
	"github.com/mailchain/mailchain/cmd/internal/datastore/os"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func substrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "substrate",
		Short:            "run substrate sequential processor",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			network, _ := cmd.Flags().GetString("network")
			protocol, _ := cmd.Flags().GetString("protocol")
			blockNumber, _ := cmd.Flags().GetUint64("start-block")

			addressRPC, _ := cmd.Flags().GetString("rpc-address")
			if addressRPC == "" {
				return errors.New("rpc-address must not be empty")
			}

			rawStorePath, _ := cmd.Flags().GetString("raw-store-path")

			connIndexer, err := newPostgresConnection(cmd, "indexer")
			if err != nil {
				return err
			}

			connPublicKey, err := newPostgresConnection(cmd, "pubkey")
			if err != nil {
				return err
			}

			connEnvelope, err := newPostgresConnection(cmd, "envelope")
			if err != nil {
				return err
			}

			defer connIndexer.Close()
			defer connPublicKey.Close()
			defer connEnvelope.Close()

			seqProcessor, err := createSubstrateProcessor(connIndexer, connPublicKey, connEnvelope, blockNumber, protocol, network, rawStorePath, addressRPC)
			if err != nil {
				return err
			}

			doSequential(seqProcessor)
		},
	}

	cmd.Flags().Uint64("start-block", 0, "Block number from which the indexer will start")
	cmd.Flags().String("protocol", protocols.Substrate, "Protocol to run against")
	cmd.Flags().String("network", substrate.EdgewareMainnet, "Network to run against")
	cmd.Flags().String("rpc-address", "", "Substrate RPC-JSON address")

	return cmd
}

func createSubstrateProcessor(connIndexer, connPublicKey, connEnvelope *sqlx.DB, blockNumber uint64, protocol, network, rawStorePath, addressRPC string) (*processor.Sequential, error) {
	ctx := context.Background()

	subClient, err := sub.NewRPC(addressRPC)
	if err != nil {
		return nil, err
	}

	syncStore, err := pq.NewSyncStore(connIndexer)
	if err != nil {
		return nil, err
	}

	pubKeyStore, err := pq.NewPublicKeyStore(connPublicKey)
	if err != nil {
		return nil, err
	}

	transactionStore, err := pq.NewTransactionStore(connEnvelope)
	if err != nil {
		return nil, err
	}

	rawStore, err := os.NewRawTransactionStore(rawStorePath)
	if err != nil {
		return nil, err
	}

	processorTransaction := sub.NewTransactionProcessor(
		transactionStore,
		rawStore,
		pubKeyStore,
	)

	if err := syncStore.PutBlockNumber(ctx, protocol, network, blockNumber); err != nil {
		return nil, err
	}

	return processor.NewSequential(
		protocols.Substrate,
		network,
		syncStore,
		sub.NewBlockProcessor(processorTransaction),
		subClient,
	), nil
}
