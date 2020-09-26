package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func rootCmd() *cobra.Command {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", ""))
	viper.SetEnvPrefix("INDEXER")
	viper.AutomaticEnv()

	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "Mailchain indexer",
	}

	cmd.AddCommand(ethereumCmd())
	cmd.AddCommand(substrateCmd())
	cmd.AddCommand(databaseCmd())

	cmd.PersistentFlags().String("postgres-host", "localhost", "Postgres server host")
	_ = viper.BindPFlag("postgres_host", cmd.PersistentFlags().Lookup("postgres-host"))
	cmd.PersistentFlags().String("postgres-sslmode", "disable", "Use SSL when connecting to Postgres")
	_ = viper.BindPFlag("postgres_sslmode", cmd.PersistentFlags().Lookup("postgres-sslmode"))
	cmd.PersistentFlags().Int("postgres-port", 5432, "Postgres server port")
	_ = viper.BindPFlag("postgres_port", cmd.PersistentFlags().Lookup("postgres-port"))

	cmd.PersistentFlags().String("indexer-postgres-user", "indexer", "Indexer postgres database user")
	_ = viper.BindPFlag("indexer_postgres_user", cmd.PersistentFlags().Lookup("indexer-postgres-user"))
	cmd.PersistentFlags().String("indexer-postgres-password", "", "Indexer postgres database password")
	_ = viper.BindPFlag("indexer_postgres_password", cmd.PersistentFlags().Lookup("indexer-postgres-password"))
	cmd.PersistentFlags().String("indexer-postgres-name", "indexer", "Indexer postgres database name")
	_ = viper.BindPFlag("indexer_postgres_name", cmd.PersistentFlags().Lookup("indexer-postgres-name"))

	cmd.PersistentFlags().String("pubkey-postgres-user", "pubkey", "Public key postgres database user")
	_ = viper.BindPFlag("pubkey_postgres_user", cmd.PersistentFlags().Lookup("pubkey-postgres-user"))
	cmd.PersistentFlags().String("pubkey-postgres-password", "", "Public key postgres database password")
	_ = viper.BindPFlag("pubkey_postgres_password", cmd.PersistentFlags().Lookup("pubkey-postgres-password"))
	cmd.PersistentFlags().String("pubkey-postgres-name", "pubkey", "Public key postgres database name")
	_ = viper.BindPFlag("pubkey_postgres_name", cmd.PersistentFlags().Lookup("pubkey-postgres-name"))

	cmd.PersistentFlags().String("envelope-postgres-user", "envelope", "Envelopes postgres database user")
	_ = viper.BindPFlag("envelope_postgres_user", cmd.PersistentFlags().Lookup("envelope-postgres-user"))
	cmd.PersistentFlags().String("envelope-postgres-password", "", "Envelopes postgres database password")
	_ = viper.BindPFlag("envelope_postgres_password", cmd.PersistentFlags().Lookup("envelope-postgres-password"))
	cmd.PersistentFlags().String("envelope-postgres-name", "envelope", "Envelopes postgres database name")
	_ = viper.BindPFlag("envelope_postgres_name", cmd.PersistentFlags().Lookup("envelope-postgres-name"))

	cmd.PersistentFlags().String("raw-store-path", "", "Path where raw transactions are stored")
	_ = viper.BindPFlag("raw_store_path", cmd.PersistentFlags().Lookup("raw-store-path"))
	cmd.PersistentFlags().Uint64("max-retries", 10, "Maximum number of retry for failures")
	_ = viper.BindPFlag("max_retries", cmd.PersistentFlags().Lookup("max-retries"))

	// Explicitly bind some flags

	return cmd
}
