package pq

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/pkg/errors"
)

// SyncStore database connection object
type SyncStore struct {
	db *sqlx.DB
}

// NewSyncStore create new postgres database
func NewSyncStore(db *sqlx.DB) (datastore.SyncStore, error) {
	return &SyncStore{db: db}, nil
}

type sync struct {
	Protocol uint8 `db:"protocol"`
	Network  uint8 `db:"network"`

	BlockNo uint64 `db:"block_no"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s SyncStore) GetBlockNumber(ctx context.Context, protocol, network string) (blockNo uint64, err error) {
	p, n, err := getProtocolNetworkUint8(protocol, network)
	if err != nil {
		return 0, errors.WithStack((err))
	}

	sql, args, err := squirrel.Select("block_no").
		From("sync").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"protocol": p}).
		Where(squirrel.Eq{"network": n}).
		ToSql()
	if err != nil {
		return 0, errors.WithStack((err))
	}

	state := sync{}
	if err := s.db.Get(&state, sql, args...); err != nil {
		return 0, errors.WithStack(err)
	}

	return state.BlockNo, nil
}

func (s SyncStore) PutBlockNumber(ctx context.Context, protocol, network string, blockNo uint64) error {
	p, n, err := getProtocolNetworkUint8(protocol, network)
	if err != nil {
		return errors.WithStack((err))
	}

	sql, args, err := squirrel.Update("sync").
		Set("block_no", blockNo).
		Set("updated_at", time.Now()).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"protocol": p}).
		Where(squirrel.Eq{"network": n}).
		ToSql()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = s.db.Exec(sql, args...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
