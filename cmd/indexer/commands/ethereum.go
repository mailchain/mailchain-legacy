package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	eth "github.com/mailchain/mailchain/cmd/indexer/internal/ethereum"
	"github.com/mailchain/mailchain/cmd/indexer/internal/processor"
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

			addressRPC, _ := cmd.Flags().GetString("rpc_address")
			if addressRPC == "" {
				return errors.New("rpc_address must not be empty")
			}

			conn, err := newPostgresConnection(cmd)
			if err != nil {
				return err
			}

			defer conn.Close()

			seqProcessor, err := createEthereumProcessor(conn, network, addressRPC)
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
	cmd.Flags().String("rpc_address", "", "Ethereum RPC-JSON address")

	return cmd
}

func createEthereumProcessor(conn *sqlx.DB, network, addressRPC string) (*processor.Sequential, error) {
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

	// TODO: rawStore does not exist yet but there is a PR https://github.com/mailchain/mailchain/pull/517
	// os.NewRawTransactionStore()

	processorTransaction := eth.NewTransactionProcessor(
		nil, // TODO: transactionStore,
		nil, // TODO: rawStore,
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
