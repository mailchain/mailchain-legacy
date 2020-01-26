package commands

import (
	"io/ioutil"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/internal/datastore/pq"
	"github.com/spf13/cobra"
)

// migration source file names
const (
	pubKeyIndexDown = "000001_public-key-index.down.sql"
	pubKeyIndexUp   = "000001_public-key-index.up.sql"

	syncUp   = "000002_sync.up.sql"
	syncDown = "000002_sync.down.sql"
)

func sslMode(useSSL bool) string {
	if useSSL {
		return "enable"
	}

	return "disable"
}

// NewPostgresConnection returns a connection to a postres database.
// The arguments are parsed from cmd.
func newPostgresConnection(cmd *cobra.Command) (*sqlx.DB, error) {
	host, _ := cmd.Flags().GetString("postgres-host")
	port, _ := cmd.Flags().GetInt("postgres-port")
	useSSL, _ := cmd.Flags().GetBool("postgres-ssl")

	user, err := cmd.Flags().GetString("postgres-user")
	if err != nil {
		return nil, err
	}

	psswd, err := cmd.Flags().GetString("postgres-password")
	if err != nil {
		return nil, err
	}

	dbname, err := cmd.Flags().GetString("postgres-name")
	if err != nil {
		return nil, err
	}

	// use default dbname is not provided
	if dbname == "" {
		dbname = user
	}

	conn, err := pq.NewConnection(user, psswd, dbname, host, sslMode(useSSL), port)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// execSQLStatements executes all SQL statements inside file.
func execSQLStatements(db *sqlx.DB, file string) error {
	for _, statement := range strings.Split(file, ";\n") {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

// execSQLFiles executes the files in given order.
func execSQLFiles(db *sqlx.DB, fileNames ...string) error {
	for _, name := range fileNames {
		file, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}

		if err := execSQLStatements(db, string(file)); err != nil {
			return err
		}
	}

	return nil
}
