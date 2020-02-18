package commands

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newMasterConnection(cmd *cobra.Command) (*sqlx.DB, error) {
	host, _ := cmd.Flags().GetString("postgres-host")
	port, _ := cmd.Flags().GetInt("postgres-port")
	sslmode, _ := cmd.Flags().GetString("postgres-sslmode")

	user, err := cmd.Flags().GetString("master-postgres-user")
	if err != nil {
		return nil, err
	}

	password, err := cmd.Flags().GetString("master-postgres-password")
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

// newPostgresConnection returns a connection to a postgres database.
// The arguments are parsed from cmd.
func newPostgresConnection(cmd *cobra.Command, kind string) (*sqlx.DB, error) {
	host, _ := cmd.Flags().GetString("postgres-host")
	port, _ := cmd.Flags().GetInt("postgres-port")
	sslmode, _ := cmd.Flags().GetString("postgres-sslmode")

	user, _ := cmd.Flags().GetString(kind + "-postgres-user")
	if user == "" {
		return nil, errors.Errorf("flag must not be empty: %s-postgres-user", kind)
	}

	password, _ := cmd.Flags().GetString(kind + "-postgres-password")
	if password == "" {
		return nil, errors.Errorf("flag must not be empty: %s-postgres-password", kind)
	}

	dbname, _ := cmd.Flags().GetString(kind + "-postgres-name")
	if dbname == "" {
		return nil, errors.Errorf("flag must not be empty: %s-postgres-name", kind)
	}

	// use default dbname, if not provided
	if dbname == "" {
		dbname = user
	}

	return pq.NewConnection(user, password, dbname, host, sslmode, port)
}
