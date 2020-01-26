package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	eth "github.com/mailchain/mailchain/cmd/indexer/internal/ethereum"
	"github.com/mailchain/mailchain/cmd/indexer/internal/processor"
	"github.com/mailchain/mailchain/cmd/internal/datastore/os"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/mailchain/mailchain/internal/protocols"
	"github.com/mailchain/mailchain/internal/protocols/ethereum"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ethereumCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "ethereum",
		Short:            "run ethereum sequential processor",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			network, _ := cmd.Flags().GetString("network")

			addressRPC, _ := cmd.Flags().GetString("rpc-address")
			if addressRPC == "" {
				return errors.New("rpc-address must not be empty")
			}

			rawStorePath, _ := cmd.Flags().GetString("raw-store-path")

			conn, err := newPostgresConnection(cmd)
			if err != nil {
				return err
			}

			defer conn.Close()

			seqProcessor, err := createEthereumProcessor(conn, network, rawStorePath, addressRPC)
			if err != nil {
				return err
			}

			for true {
				if err := seqProcessor.NextBlock(context.Background()); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%+v", err)
				}
				fmt.Println("Infinite Loop 2")
				time.Sleep(time.Second)
			}

			return nil
		},
	}
	cmd.Flags().String("network", ethereum.Mainnet, "Network to run against")
	cmd.Flags().String("rpc-address", "", "Ethereum RPC-JSON address")

	return cmd
}

func createEthereumProcessor(conn *sqlx.DB, network, rawStorePath, addressRPC string) (*processor.Sequential, error) {
	ethClient, err := eth.NewRPC(addressRPC)
	if err != nil {
		return nil, err
	}

	networkID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	syncStore, err := pq.NewSyncStore(conn)
	if err != nil {
		return nil, err
	}

	pubKeyStore, err := pq.NewPublicKeyStore(conn)
	if err != nil {
		return nil, err
	}

	// TODO: transaction store does not exist yet issue/515
	// transactionStore, err := pq.NewTransactionStore(conn)
	// if err != nil {
	// 	return err
	// }

	rawStore, err := os.NewRawTransactionStore(rawStorePath)
	if err != nil {
		return nil, err
	}

	processorTransaction := eth.NewTransactionProcessor(
		nil, // TODO: transactionStore,
		rawStore,
		pubKeyStore,
		networkID,
	)

	return processor.NewSequential(
		protocols.Ethereum,
		network,
		syncStore,
		eth.NewBlockProcessor(processorTransaction),
		ethClient,
	), nil
}
