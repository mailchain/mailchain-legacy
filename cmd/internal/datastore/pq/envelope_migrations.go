package pq

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

func MigrateEnvelope(db *sql.DB, up bool) (int, error) {
	direction := migrate.Down
	if up {
		direction = migrate.Up
	}

	return migrate.Exec(db, "postgres",
		&migrate.MemoryMigrationSource{
			Migrations: []*migrate.Migration{
				{Id: "1581972758197-create-transactions-table",
					Up: []string{`
					CREATE TABLE IF NOT EXISTS transactions(
						-- Primary Key
						protocol        SMALLINT NOT NULL,
						network         SMALLINT NOT NULL,
						hash            BYTEA NOT NULL,
						-- Values
						tx_from            BYTEA NOT NULL,
						tx_to              BYTEA NOT NULL,
						tx_data            BYTEA NOT NULL,
						tx_block_no        BIGINT NOT NULL,
						tx_block_hash      BYTEA NOT NULL,    
						tx_value           BYTEA NOT NULL,
						tx_gas_used        BYTEA NOT NULL,
						tx_gas_price       BYTEA NOT NULL,
						PRIMARY KEY(protocol, network, hash)
					);`},
					Down: []string{`DROP TABLE transactions;`},
				},
			},
		},
		direction,
	)
}
