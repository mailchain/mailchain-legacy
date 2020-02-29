package commands

import (
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "Mailchain indexer",
	}

	cmd.AddCommand(ethereumCmd())
	cmd.AddCommand(databaseCmd())

	cmd.PersistentFlags().String("postgres-host", "localhost", "Postgres server host")
	cmd.PersistentFlags().String("postgres-sslmode", "disable", "Use SSL when connecting to Postgres")
	cmd.PersistentFlags().Int("postgres-port", 5432, "Postgres server port")

	cmd.PersistentFlags().String("indexer-postgres-user", "indexer", "Indexer postgres database user")
	cmd.PersistentFlags().String("indexer-postgres-password", "", "Indexer postgres database password")
	cmd.PersistentFlags().String("indexer-postgres-name", "indexer", "Indexer postgres database name")

	cmd.PersistentFlags().String("pubkey-postgres-user", "pubkey", "Public key postgres database user")
	cmd.PersistentFlags().String("pubkey-postgres-password", "", "Public key postgres database password")
	cmd.PersistentFlags().String("pubkey-postgres-name", "pubkey", "Public key postgres database name")

	cmd.PersistentFlags().String("envelope-postgres-user", "envelope", "Envelopes postgres database user")
	cmd.PersistentFlags().String("envelope-postgres-password", "", "Envelopes postgres database password")
	cmd.PersistentFlags().String("envelope-postgres-name", "envelope", "Envelopes postgres database name")

	cmd.PersistentFlags().String("raw-store-path", "", "Path where raw transactions are stored")

	return cmd
}
