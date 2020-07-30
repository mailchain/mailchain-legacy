package commands

import (
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
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

	db, err := pq.NewConnection(user, password, "envelope", host, sslmode, port)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open connection: %s", host)
	}

	return db, nil
}
