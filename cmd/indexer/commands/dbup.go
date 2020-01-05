package commands

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
)

func dbUpCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:              "up",
		Short:            "init postgres database for mailchain indexer",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newPostgresConnection(cmd)
			if err != nil {
				return err
			}

			defer conn.Close()

			p, err := cmd.Flags().GetString("path")
			if err != nil {
				return err
			}

			files := []string{
				path.Join(p, pubKeyIndexUp),
				path.Join(p, syncUp),
			}

			if err = execSQLFiles(conn, files...); err != nil {
				return err
			}

			fmt.Println("successfully initialized")
			return nil
		},
	}

	return cmd, nil
}
