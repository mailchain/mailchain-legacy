package pq

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // needed for postgres to work
	"github.com/pkg/errors"
)

func NewDatabase(user, password, databaseName, host, sslmode string, port int) (*sqlx.DB, error) {
	// The first argument corresponds to the driver name that the driver
	// (in this case, `lib/pq`) used to register itself in `database/sql`.
	// The next argument specifies the parameters to be used in the connection.
	// Details about this string can be seen at https://godoc.org/github.com/lib/pq
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		user, password, databaseName, host, port, sslmode))
	if err != nil {
		return nil, errors.Wrapf(err,
			"Couldn't open connection to postgre database (%s)", databaseName)
	}

	// Ping verifies if the connection to the database is alive or if a
	// new connection can be made.
	if err = db.Ping(); err != nil {
		return nil, errors.Wrapf(err,
			"Couldn't ping postgre database (%s)", databaseName)
	}

	return db, nil
}
