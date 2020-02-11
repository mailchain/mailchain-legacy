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

	cmd.PersistentFlags().String("postgres-host", "localhost", "Postgres server host")
	cmd.PersistentFlags().String("postgres-user", "", "Postgres database user")
	cmd.PersistentFlags().String("postgres-password", "", "Postgres database password")
	cmd.PersistentFlags().String("postgres-name", "", "Postgres database name")
	cmd.PersistentFlags().Bool("postgres-ssl", false, "Use SSL when connecting to Postgres")
	cmd.PersistentFlags().Int("postgres-port", 5432, "Postgres server port")

	cmd.PersistentFlags().String("raw-store-path", "", "Path where raw transactions are stored")

	return cmd
}
