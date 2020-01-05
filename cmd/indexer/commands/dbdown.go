package commands

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
)

func dbDownCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:              "down",
		Short:            "destroy postgres database for mailchain indexer",
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
				path.Join(p, pubKeyIndexDown),
				path.Join(p, syncDown),
			}

			if err = execSQLFiles(conn, files...); err != nil {
				return err
			}

			fmt.Println("successfully destroyed")
			return nil
		},
	}

	return cmd, nil
}
