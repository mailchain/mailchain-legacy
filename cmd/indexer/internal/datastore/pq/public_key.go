package pq

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/indexer/internal/datastore"
	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/multikey"
	"github.com/pkg/errors"
)

// PublicKeyStore database connection object
type PublicKeyStore struct {
	db *sqlx.DB
}

// PublicKeyStore create new pointer to postgres database
func NewPublicKeyStore(db *sqlx.DB) (datastore.PublicKeyStore, error) {
	return &PublicKeyStore{db: db}, nil
}

type public_key struct {
	Protocol uint8  `db:"protocol"`
	Network  uint8  `db:"network"`
	Address  []byte `db:"address"`

	PublicKeyType uint8  `db:"public_key_type"`
	PublicKey     []byte `db:"public_key"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s PublicKeyStore) PutPublicKey(ctx context.Context, protocol, network string, address []byte, pubKey crypto.PublicKey) error {
	p, n, err := getProtocolNetworkUint8(protocol, network)
	if err != nil {
		return errors.WithStack((err))
	}

	uPublicKeyType, err := getPublicKeyTypeUint8(pubKey.Kind())
	if err != nil {
		return errors.WithStack((err))
	}

	sql, args, err := squirrel.Update("public_keys").
		Set("public_key_type", uPublicKeyType).
		Set("public_key", pubKey.Bytes()).
		Set("updated_at", time.Now()).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"protocol": p}).
		Where(squirrel.Eq{"network": n}).
		Where(squirrel.Eq{"address": address}).
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

func (s PublicKeyStore) GetPublicKey(ctx context.Context, protocol, network string, address []byte) (pubKey crypto.PublicKey, err error) {
	p, n, err := getProtocolNetworkUint8(protocol, network)
	if err != nil {
		return nil, errors.WithStack((err))
	}

	sql, args, err := squirrel.Select("public_key_type", "public_key").
		From("public_keys").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"protocol": p}).
		Where(squirrel.Eq{"network": n}).
		Where(squirrel.Eq{"address": address}).
		ToSql()
	if err != nil {
		return nil, errors.WithStack((err))
	}

	state := public_key{}
	if err := s.db.Get(&state, sql, args...); err != nil {
		return nil, errors.WithStack(err)
	}

	publicKeyType, err := getPublicKeyTypeString(state.PublicKeyType)
	if err != nil {
		return nil, errors.WithStack((err))
	}

	publicKey, err := multikey.PublicKeyFromBytes(publicKeyType, state.PublicKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return publicKey, nil
}
