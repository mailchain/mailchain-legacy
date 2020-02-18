package pq

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

func MigratePublicKey(db *sql.DB, up bool) (int, error) {
	direction := migrate.Down
	if up {
		direction = migrate.Up
	}

	return migrate.Exec(db, "postgres",
		&migrate.MemoryMigrationSource{
			Migrations: []*migrate.Migration{
				{Id: "1581972558643-create-public-key-table",
					Up: []string{`
					CREATE TABLE IF NOT EXISTS public_keys(
						-- Primary Key    
						protocol                SMALLINT NOT NULL,
						network                 SMALLINT NOT NULL,
						address                 BYTEA NOT NULL,
						-- Values
						public_key_type         SMALLINT NOT NULL,
						public_key              BYTEA NOT NULL,
						created_block_hash      BYTEA NOT NULL,
						updated_block_hash      BYTEA NOT NULL,
						created_tx_hash         BYTEA NOT NULL,
						updated_tx_hash         BYTEA NOT NULL,
						PRIMARY KEY(protocol, network, address)
					);`},
					Down: []string{`DROP TABLE public_keys;`},
				},
			},
		},
		direction,
	)
}
