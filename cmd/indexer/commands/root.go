package commands

import (
	"github.com/spf13/cobra"
)

func rootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "indexer",
		Short: "Mailchain indexer",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newPostgresConnection(cmd)
			if err != nil {
				return err
			}

			defer conn.Close()

			panic("implement me!")
		},
	}

	dbInit, err := dbUpCmd()
	if err != nil {
		return nil, err
	}

	dbDestroy, err := dbDownCmd()
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(dbInit)
	cmd.AddCommand(dbDestroy)

	cmd.PersistentFlags().String("host", "localhost", "host url")
	cmd.PersistentFlags().Int("port", 5432, "")
	cmd.PersistentFlags().String("user", "", "")
	cmd.PersistentFlags().String("password", "", "")
	cmd.PersistentFlags().String("dbname", "", "")
	cmd.PersistentFlags().Bool("ssl", false, "")
	cmd.PersistentFlags().String("path", "", "path to migration source files")

	return cmd, nil
}
