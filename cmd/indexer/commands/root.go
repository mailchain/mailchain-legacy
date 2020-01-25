package commands

import (
	"github.com/spf13/cobra"
)

func rootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "Mailchain indexer",
	}

	dbInit, err := dbUpCmd()
	if err != nil {
		return nil, err
	}

	dbDestroy, err := dbDownCmd()
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(ethereumCmd())
	cmd.AddCommand(dbInit)
	cmd.AddCommand(dbDestroy)

	cmd.PersistentFlags().String("postgres_host", "localhost", "Postgres server host")
	cmd.PersistentFlags().Int("postgres_port", 5432, "Postgres server port")
	cmd.PersistentFlags().String("postgres_user", "", "Postgres database user")
	cmd.PersistentFlags().String("postgres_password", "", "Postgres database password")
	cmd.PersistentFlags().String("postgres_name", "", "Postgres database name")
	cmd.PersistentFlags().Bool("postgres_ssl", false, "Use SSL when connecting to Postgres")
	cmd.PersistentFlags().String("path", "", "path to migration source files")

	return cmd, nil
}
