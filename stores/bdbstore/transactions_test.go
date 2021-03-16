package bdbstore

import (
	"path"
	"sort"
	"testing"

	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/addressing/addressingtest"
	"github.com/mailchain/mailchain/stores"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_PutTransaction(t *testing.T) {
	type args struct {
		protocol string
		network  string
		address  []byte
		tx       stores.Transaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"success",
			args{
				"ethereum",
				"mainnet",
				addressingtest.EthereumCharlotte,
				stores.Transaction{
					EnvelopeData: []byte("env1"),
					BlockNumber:  100,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown, _ := setupDB(path.Join(getTempDir(), tt.name))
			defer teardown()

			if err := db.PutTransaction(tt.args.protocol, tt.args.network, tt.args.address, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("Database.PutTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_GetTransactions(t *testing.T) {
	type args struct {
		protocol string
		network  string
		address  []byte
		tx       stores.Transaction
	}
	tests := []struct {
		name          string
		args          []args
		queryProtocol string
		queryNetwork  string
		queryAddress  []byte
		want          []stores.Transaction
		wantErr       bool
	}{
		{
			"single-tx-charlotte",
			[]args{
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumCharlotte,
					stores.Transaction{
						EnvelopeData: []byte("env2"),
						BlockNumber:  2,
					},
				},
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumSofia,
					stores.Transaction{
						EnvelopeData: []byte("env3"),
						BlockNumber:  3,
					},
				},
			},
			"ethereum",
			"mainnet",
			addressingtest.EthereumCharlotte,
			[]stores.Transaction{
				{
					EnvelopeData: []byte("env2"),
					BlockNumber:  2,
				},
			},
			false,
		},
		{
			"multiple-message-charlotte",
			[]args{
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumCharlotte,
					stores.Transaction{
						EnvelopeData: []byte("env2"),
						BlockNumber:  2,
					},
				},
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumSofia,
					stores.Transaction{
						EnvelopeData: []byte("env3"),
						BlockNumber:  3,
					},
				},
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumCharlotte,
					stores.Transaction{
						EnvelopeData: []byte("env4"),
						BlockNumber:  4,
					},
				},
			},
			"ethereum",
			"mainnet",
			addressingtest.EthereumCharlotte,
			[]stores.Transaction{
				{
					EnvelopeData: []byte("env4"),
					BlockNumber:  4,
				},
				{
					EnvelopeData: []byte("env2"),
					BlockNumber:  2,
				},
			},
			false,
		},
		{
			"multiple-message-charlotte-irregular-entry-order",
			[]args{
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumCharlotte,
					stores.Transaction{
						EnvelopeData: []byte("env2"),
						BlockNumber:  2,
					},
				},
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumSofia,
					stores.Transaction{
						EnvelopeData: []byte("env2"),
						BlockNumber:  3,
					},
				},
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumCharlotte,
					stores.Transaction{
						EnvelopeData: []byte("env4"),
						BlockNumber:  4,
					},
				},
				{
					"ethereum",
					"mainnet",
					addressingtest.EthereumCharlotte,
					stores.Transaction{
						EnvelopeData: []byte("env1"),
						BlockNumber:  1,
					},
				},
			},
			"ethereum",
			"mainnet",
			addressingtest.EthereumCharlotte,
			[]stores.Transaction{
				{
					EnvelopeData: []byte("env4"),
					BlockNumber:  4,
				},
				{
					EnvelopeData: []byte("env2"),
					BlockNumber:  2,
				},
				{
					EnvelopeData: []byte("env1"),
					BlockNumber:  1,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown, _ := setupDB(path.Join(getTempDir(), tt.name))
			defer teardown()

			for _, args := range tt.args {
				if err := db.PutTransaction(args.protocol, args.network, args.address, args.tx); err != nil {
					t.Errorf("PutTransaction() returned an error %v", err)
				}
			}

			got, err := db.GetTransactions(tt.queryProtocol, tt.queryNetwork, tt.queryAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactions() err = %v, want err %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(t, got, tt.want) {
				t.Errorf("Database.GetTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transactionKeyOrder(t *testing.T) {
	orderedBlockNumbers := []int64{
		999999999999,
		999999999998,
		88888888888,
		88888888887,
		7777777777,
		7777777776,
		55556,
		55555,
		4445,
		4444,
		334,
		333,
		22,
		11,
		2,
		1,
	}

	txKeysBlockNumbers := map[string]int64{}
	transactionKeys := []string{}
	for _, x := range orderedBlockNumbers {
		txKey := encoding.EncodeHex(transactionKey("0", x, nil))
		// blockNumbersTxKeys[x] = txKey
		txKeysBlockNumbers[txKey] = x
		transactionKeys = append(transactionKeys, txKey)
	}

	sort.Strings(transactionKeys)
	for i, x := range transactionKeys {
		assert.Equal(t, orderedBlockNumbers[i], txKeysBlockNumbers[x])
	}
}
