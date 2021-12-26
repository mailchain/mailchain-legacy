package pq

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mailchain/mailchain/cmd/internal/datastore"
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

type publicKey struct {
	Address          []byte `db:"address"`
	PublicKey        []byte `db:"public_key"`
	CreatedBlockHash []byte `db:"created_block_hash"`
	UpdatedBlockHash []byte `db:"updated_block_hash"`
	CreatedTxHash    []byte `db:"created_tx_hash"`
	UpdatedTxHash    []byte `db:"updated_tx_hash"`
	Protocol         uint8  `db:"protocol"`
	Network          uint8  `db:"network"`
	PublicKeyType    uint8  `db:"public_key_type"`
}

func (s PublicKeyStore) PutPublicKey(ctx context.Context, protocol, network string, address []byte, pubKey *datastore.PublicKey) error {
	p, n, err := getProtocolNetworkUint8(protocol, network)
	if err != nil {
		return errors.WithStack((err))
	}

	kind, err := multikey.KindFromPublicKey(pubKey.PublicKey)
	if err != nil {
		return errors.WithStack((err))
	}

	uPublicKeyType, err := getPublicKeyTypeUint8(kind)
	if err != nil {
		return errors.WithStack((err))
	}

	sql, args, err := squirrel.Insert("public_keys").
		Columns("protocol", "network", "address", "public_key_type", "public_key", "created_block_hash", "updated_block_hash", "created_tx_hash", "updated_tx_hash").
		Values(p, n, address, uPublicKeyType, pubKey.PublicKey.Bytes(), pubKey.BlockHash, pubKey.BlockHash, pubKey.TxHash, pubKey.TxHash).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("ON CONFLICT (protocol, network, address) DO UPDATE SET public_key_type = $10, public_key = $11, updated_block_hash = $12, updated_tx_hash = $13",
			uPublicKeyType, pubKey.PublicKey.Bytes(), pubKey.BlockHash, pubKey.TxHash).
		ToSql()
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err = s.db.Exec(sql, args...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s PublicKeyStore) GetPublicKey(ctx context.Context, protocol, network string, address []byte) (pubKey *datastore.PublicKey, err error) {
	p, n, err := getProtocolNetworkUint8(protocol, network)
	if err != nil {
		return nil, errors.WithStack((err))
	}

	sql, args, err := squirrel.Select("public_key_type", "public_key", "updated_block_hash", "updated_tx_hash").
		From("public_keys").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"protocol": p}).
		Where(squirrel.Eq{"network": n}).
		Where(squirrel.Eq{"address": address}).
		ToSql()
	if err != nil {
		return nil, errors.WithStack((err))
	}

	state := publicKey{}
	if err := s.db.Get(&state, sql, args...); err != nil {
		return nil, errors.WithStack(err)
	}

	publicKeyType, err := getPublicKeyTypeString(state.PublicKeyType)
	if err != nil {
		return nil, errors.WithStack((err))
	}

	cryptoPublicKey, err := multikey.PublicKeyFromBytes(publicKeyType, state.PublicKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	publicKey := &datastore.PublicKey{PublicKey: cryptoPublicKey, BlockHash: state.UpdatedBlockHash, TxHash: state.UpdatedTxHash}

	return publicKey, nil
}
