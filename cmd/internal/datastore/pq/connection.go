package pq

import (
	"fmt"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // needed for postgres to work
	"github.com/pkg/errors"
)

func NewConnection(user, password, databaseName, host, sslmode string, port int) (*sqlx.DB, error) {
	// The first argument corresponds to the driver name that the driver
	// (in this case, `lib/pq`) used to register itself in `database/sql`.
	// The next argument specifies the parameters to be used in the connection.
	// Details about this string can be seen at https://godoc.org/github.com/lib/pq
	var db *sqlx.DB

	operation := func() error {
		var err error
		db, err = sqlx.Connect("postgres", fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
			user, password, databaseName, host, port, sslmode))

		return err
	}

	err := backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5))

	return db, errors.Wrapf(err, "could not open connection to postgres database: %s", databaseName)
}
