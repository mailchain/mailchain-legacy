package commands

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/indexer/internal/processor"
	sub "github.com/mailchain/mailchain/cmd/indexer/internal/substrate"
	"github.com/mailchain/mailchain/cmd/internal/datastore/noop"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/substrate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func substrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "substrate",
		Short:            "run substrate sequential processor",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			network := viper.GetString("network")
			protocol := viper.GetString("protocol")
			blockNumber := viper.GetString("start_block")
			maxRetries := viper.GetUint64("max_retries")
			addressRPC := viper.GetString("rpc_address")
			if addressRPC == "" {
				return errors.Errorf("rpc-address must not be empty")
			}

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

			var subClient *sub.BlockClient
			operation := func() error {
				var err error
				subClient, err = sub.NewRPC(addressRPC)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%+v\n", err)
				}

				return err
			}

			if err := backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Number of retries has reached to %d. Exiting.\\n", 5)
				return err
			}

			seqProcessor, err := createSubstrateProcessor(subClient, connIndexer, connPublicKey, connEnvelope, blockNumber, protocol, network)
			if err != nil {
				return err
			}

			doSequential(cmd, seqProcessor, maxRetries)

			return nil
		},
	}

	cmd.Flags().String("start-block", "latest", "Block number from which the indexer will start, e.g. 10000, or 'latest'")
	_ = viper.BindPFlag("start_block", cmd.Flags().Lookup("start-block"))
	cmd.Flags().String("protocol", protocols.Substrate, "Protocol to run against")
	_ = viper.BindPFlag("protocol", cmd.Flags().Lookup("protocol"))
	cmd.Flags().String("network", substrate.EdgewareMainnet, "Network to run against")
	_ = viper.BindPFlag("network", cmd.Flags().Lookup("network"))
	cmd.Flags().String("rpc-address", "", "Substrate RPC-JSON address")
	_ = viper.BindPFlag("rpc_address", cmd.Flags().Lookup("rpc-address"))

	return cmd
}

func createSubstrateProcessor(client *sub.BlockClient, connIndexer, connPublicKey, connEnvelope *sqlx.DB, blockNumber, protocol, network string) (*processor.Sequential, error) {
	ctx := context.Background()

	blockNo, err := getBlockNumber(blockNumber, client)
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

	rawStore, err := noop.NewRawTransactionStore()
	if err != nil {
		return nil, err
	}

	processorTransaction := sub.NewExtrinsicProcessor(transactionStore, rawStore, pubKeyStore)

	if err := syncStore.PutBlockNumber(ctx, protocol, network, blockNo); err != nil {
		return nil, err
	}

	return processor.NewSequential(
		protocols.Substrate,
		network,
		syncStore,
		sub.NewBlockProcessor(processorTransaction),
		client,
	), nil
}
