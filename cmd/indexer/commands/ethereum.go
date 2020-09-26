package commands

import (
	"context"

	"github.com/ethereum/go-ethereum/params"
	"github.com/jmoiron/sqlx"
	eth "github.com/mailchain/mailchain/cmd/indexer/internal/ethereum"
	"github.com/mailchain/mailchain/cmd/indexer/internal/processor"
	"github.com/mailchain/mailchain/cmd/internal/datastore/os"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ethereumCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "ethereum",
		Short:            "run ethereum sequential processor",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			network := viper.GetString("network")
			protocol := viper.GetString("protocol")
			blockNumber := viper.GetString("start_block")

			addressRPC := viper.GetString("rpc_address")
			if addressRPC == "" {
				return errors.Errorf("rpc-address must not be empty")
			}

			rawStorePath, _ := cmd.Flags().GetString("raw-store-path")

			maxRetries, _ := cmd.PersistentFlags().GetUint64("max-retries")

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

			seqProcessor, err := createEthereumProcessor(connIndexer, connPublicKey, connEnvelope, blockNumber, protocol, network, rawStorePath, addressRPC)
			if err != nil {
				return err
			}

			doSequential(cmd, seqProcessor, maxRetries)

			return nil
		},
	}

	cmd.Flags().String("start-block", "latest", "Block number from which the indexer will start, e.g. 10000, or 'latest'")
	_ = viper.BindPFlag("start_block", cmd.Flags().Lookup("start-block"))
	cmd.Flags().String("protocol", protocols.Ethereum, "Protocol to run against")
	_ = viper.BindPFlag("protocol", cmd.Flags().Lookup("protocol"))
	cmd.Flags().String("network", ethereum.Mainnet, "Network to run against")
	_ = viper.BindPFlag("network", cmd.Flags().Lookup("network"))
	cmd.Flags().String("rpc-address", "", "Ethereum RPC-JSON address")
	_ = viper.BindPFlag("rpc_address", cmd.Flags().Lookup("rpc-address"))

	return cmd
}

func createEthereumProcessor(connIndexer, connPublicKey, connEnvelope *sqlx.DB, blockNumber, protocol, network, rawStorePath, addressRPC string) (*processor.Sequential, error) {
	ctx := context.Background()

	ethClient, err := eth.NewRPC(addressRPC)
	if err != nil {
		return nil, err
	}

	blockNo, err := getBlockNumber(blockNumber, ethClient)
	if err != nil {
		return nil, err
	}

	networkID, err := ethClient.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	chainCfg, err := chainConfig(network)
	if err != nil {
		return nil, err
	}

	if chainCfg.ChainID.Cmp(networkID) != 0 {
		return nil, errors.Errorf("networkID from RPC does not match chain config network ID")
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

	processorTransaction := eth.NewTransactionProcessor(
		transactionStore,
		rawStore,
		pubKeyStore,
		chainCfg,
	)

	if err := syncStore.PutBlockNumber(ctx, protocol, network, blockNo); err != nil {
		return nil, err
	}

	return processor.NewSequential(
		protocols.Ethereum,
		network,
		syncStore,
		eth.NewBlockProcessor(processorTransaction),
		ethClient,
	), nil
}

func chainConfig(network string) (*params.ChainConfig, error) {
	switch network {
	case ethereum.Goerli:
		return params.GoerliChainConfig, nil
	case ethereum.Mainnet:
		return params.MainnetChainConfig, nil
	case ethereum.Rinkeby:
		return params.RinkebyChainConfig, nil
	default:
		return nil, errors.Errorf("can not determine chain config from network: %s", network)
	}
}
