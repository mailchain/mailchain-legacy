package commands

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newDatabaseConnection(cmd *cobra.Command) (*sqlx.DB, error) {
	host, _ := cmd.Flags().GetString("postgres-host")
	port, _ := cmd.Flags().GetInt("postgres-port")
	sslmode, _ := cmd.Flags().GetString("postgres-sslmode")

	user, err := cmd.Flags().GetString("postgres-user")
	if err != nil {
		return nil, err
	}

	password, err := cmd.Flags().GetString("postgres-password")
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"user=%s password=%s host=%s port=%d sslmode=%s",
		user, password, host, port, sslmode))
	if err != nil {
		return nil, errors.Wrapf(err, "could not open connection: %s", host)
	}

	return db, nil
}
