package pq

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/indexer/internal/datastore"
	"github.com/pkg/errors"
)

// SyncStore database connection object
type SyncStore struct {
	db  *sqlx.DB
	now func() time.Time
}

// NewSyncStore create new postgres database
func NewSyncStore(db *sqlx.DB, now func() time.Time) (datastore.SyncStore, error) {
	return &SyncStore{db: db, now: now}, nil
}

type sync struct {
	Protocol uint8 `db:"protocol"`
	Network  uint8 `db:"network"`

	BlockNo uint64 `db:"block_no"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s SyncStore) GetBlockNumber(ctx context.Context, protocol, network string) (blockNo uint64, err error) {
	uProtocol, ok := protocolUint8[protocol]
	if !ok {
		return 0, errors.Errorf("unknown protocol: %q", protocol)
	}

	uNetwork, ok := protocolNetworkUint8[protocol][network]
	if !ok {
		return 0, errors.Errorf("unknown protocol.network: \"%s.%s\"", protocol, network)
	}

	sql, args, err := s.selectBlockNumberQuery(uProtocol, uNetwork)
	if err != nil {
		return 0, err
	}

	//  // You can also get a single result, a la QueryRow
	state := sync{}

	err = s.db.Get(&state, sql, args)
	if err != nil {
		return 0, err
	}

	return state.BlockNo, nil
}

func (s SyncStore) PutBlockNumber(ctx context.Context, protocol, network string, blockNo uint64) error {
	uProtocol, ok := protocolUint8[protocol]
	if !ok {
		return errors.Errorf("unknown protocol: %q", protocol)
	}

	uNetwork, ok := protocolNetworkUint8[protocol][network]
	if !ok {
		return errors.Errorf("unknown protocol.network: \"%s.%s\"", protocol, network)
	}

	sql, args, err := s.updateBlockNumberQuery(uProtocol, uNetwork, blockNo)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = s.db.Exec(sql, args...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s SyncStore) updateBlockNumberQuery(protocol, network uint8, blockNo uint64) (sql string, args []interface{}, err error) {
	return squirrel.Update("sync").
		Set("block_no", blockNo).
		Set("updated_at", s.now()).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"protocol": protocol}).
		Where(squirrel.Eq{"network": network}).
		ToSql()
}

func (s SyncStore) selectBlockNumberQuery(protocol, network uint8) (sql string, args []interface{}, err error) {
	return squirrel.Select("block_no").
		From("sync").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"protocol": protocol}).
		Where(squirrel.Eq{"network": network}).
		ToSql()
}
