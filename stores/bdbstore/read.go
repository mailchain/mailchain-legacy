package bdbstore

import (
	"fmt"

	"github.com/dgraph-io/badger/v2"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/pkg/errors"
)

// PutMessageRead mark message as read.
func (db *Database) PutMessageRead(messageID mail.ID) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Set(db.messageReadKey(messageID), []byte{1})
	})
}

// DeleteMessageRead mark message as unread.
func (db *Database) DeleteMessageRead(messageID mail.ID) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(db.messageReadKey(messageID))
	})
}

// GetReadStatus returns true if message is read.
func (db *Database) GetReadStatus(messageID mail.ID) (bool, error) {
	var (
		val []byte
		err error
	)

	err = db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(db.messageReadKey(messageID))
		if err != nil {
			return err
		}

		// copy the value as the item is only valid
		// while the txn is open.
		val, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		return false, err
	}

	if len(val) != 1 {
		return false, errors.Errorf("invalid read status length")
	}

	return val[0] == 1, nil
}

func (db *Database) messageReadKey(messageID mail.ID) []byte {
	return []byte(fmt.Sprintf("message.%s.read", messageID.HexString()))
}
