package commands

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	libpq "github.com/lib/pq"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func databaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "configure database",
	}

	cmd.AddCommand(upCmd())
	cmd.AddCommand(downCmd())

	cmd.PersistentFlags().String("master-postgres-user", "", "Postgres database user")
	cmd.PersistentFlags().String("master-postgres-password", "", "Postgres database password")

	return cmd
}

func upCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "bring database up to latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			masterConn, err := newMasterConnection(cmd)
			if err != nil {
				return err
			}

			if err := upsertDatabase(cmd, masterConn, "indexer"); err != nil {
				return errors.WithStack(err)
			}

			if err := migrateDatabase(cmd, "indexer", true, pq.MigrateIndexer); err != nil {
				return errors.WithStack(err)
			}

			if err := upsertDatabase(cmd, masterConn, "pubkey"); err != nil {
				return errors.WithStack(err)
			}

			if err := migrateDatabase(cmd, "pubkey", true, pq.MigratePublicKey); err != nil {
				return errors.WithStack(err)
			}

			if err := upsertDatabase(cmd, masterConn, "envelope"); err != nil {
				return errors.WithStack(err)
			}

			if err := migrateDatabase(cmd, "envelope", true, pq.MigrateEnvelope); err != nil {
				return errors.WithStack(err)
			}

			fmt.Printf("competed!\n")

			return nil
		},
	}
}

func downCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "migrate down the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := migrateDatabase(cmd, "indexer", false, pq.MigrateIndexer); err != nil {
				return errors.WithStack(err)
			}

			if err := migrateDatabase(cmd, "pubkey", false, pq.MigratePublicKey); err != nil {
				return errors.WithStack(err)
			}

			if err := migrateDatabase(cmd, "envelope", false, pq.MigrateEnvelope); err != nil {
				return errors.WithStack(err)
			}

			fmt.Printf("competed!\n")

			return nil
		},
	}
}

func upsertDatabase(cmd *cobra.Command, connMaster *sqlx.DB, kind string) error {
	var err error

	password, _ := cmd.Flags().GetString(kind + "-postgres-password")
	if password == "" {
		return errors.Errorf("must not be empty: %s-postgres-password", kind)
	}

	user, _ := cmd.Flags().GetString(kind + "-postgres-user")
	name, _ := cmd.Flags().GetString(kind + "-postgres-name")

	_, err = connMaster.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", user, password))
	if err != nil {
		if pqErr, ok := err.(*libpq.Error); !ok || pqErr.Code != "42710" {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "skipping create user, already exists user: %s\n", user)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "user created: %s\n", user)
	}

	_, err = connMaster.Exec(fmt.Sprintf("CREATE DATABASE %s;", name))
	if err != nil {
		if pqErr, ok := err.(*libpq.Error); !ok || pqErr.Code != "42P04" {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "skipping create database, already exists database: %s\n", name)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "database created: %s\n", name)
	}

	_, err = connMaster.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;", name, user))
	if err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "skipping grant privileges, already exists on database: %s\n", name)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "privileges granted on database: %s\n", name)
	}

	return nil
}

func migrateDatabase(cmd *cobra.Command, name string, up bool, migrateFunc func(db *sql.DB, up bool) (int, error)) error {
	conn, err := newPostgresConnection(cmd, name)
	if err != nil {
		return errors.WithStack(err)
	}

	defer conn.Close()

	n, err := migrateFunc(conn.DB, up)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "migrated database: %s\n", name)
	fmt.Fprintf(cmd.OutOrStdout(), "applied migration files: %d\n", n)

	return nil
}
