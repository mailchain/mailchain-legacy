package pq

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

func MigrateIndexer(db *sql.DB, up bool) (int, error) {
	direction := migrate.Down
	if up {
		direction = migrate.Up
	}

	return migrate.Exec(db, "postgres",
		&migrate.MemoryMigrationSource{
			Migrations: []*migrate.Migration{
				{Id: "1581972582040-create-sync-table",
					Up: []string{`
					CREATE TABLE IF NOT EXISTS sync(
						-- Primary Key    
						protocol                SMALLINT NOT NULL,
						network                 SMALLINT NOT NULL,
						-- Values
						block_no                BIGINT NOT NULL,
						-- Metadata
						created_at              TIMESTAMP NOT NULL,
						updated_at              TIMESTAMP NOT NULL,
						PRIMARY KEY(protocol, network)
					);`},
					Down: []string{`DROP TABLE sync;`},
				},
			},
		},
		direction,
	)
}
