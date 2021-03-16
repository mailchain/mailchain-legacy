package bdbstore

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/badger/v2"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/stores"
	"github.com/multiformats/go-multihash"
)

func (db *Database) PutTransaction(protocol, network string, address []byte, tx stores.Transaction) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(tx); err != nil {
		return err
	}

	enc := buf.Bytes()

	return db.db.Update(func(txn *badger.Txn) error {
		prefixKey := getTransactionPrefixKey(protocol, network, address)

		key := transactionKey(prefixKey, tx.BlockNumber, enc)
		return txn.Set(key, enc)
	})
}

func (db *Database) GetTransactions(protocol, network string, address []byte) ([]stores.Transaction, error) {
	var txs []stores.Transaction

	err := db.db.View(func(txn *badger.Txn) error {
		prefix := []byte(getTransactionPrefixKey(protocol, network, address))
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var tx stores.Transaction
			if err := gob.NewDecoder(bytes.NewReader(val)).Decode(&tx); err != nil {
				return err
			}

			txs = append(txs, tx)
		}
		return nil
	})

	return txs, err
}

func transactionKey(prefixKey string, order int64, encodedTx []byte) []byte {
	id, _ := multihash.Sum(encodedTx, multihash.SHA3_256, -1)

	orderBytes := new(bytes.Buffer)
	_ = binary.Write(orderBytes, binary.BigEndian, order*-1)

	return []byte(fmt.Sprintf("%s/%s/%s", prefixKey, encoding.EncodeHexZeroX(orderBytes.Bytes()), encoding.EncodeHexZeroX(id)))
}

func getTransactionPrefixKey(protocol, network string, address []byte) string {
	return fmt.Sprintf("%s.%s.%s", protocol, network, encoding.EncodeHex(address))
}
