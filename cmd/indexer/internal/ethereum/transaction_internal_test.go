package ethereum

import (
	"math/big"
	"testing"

	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/datastore/datastoretest"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionProcessor(t *testing.T) {
	type args struct {
		store     datastore.TransactionStore
		rawStore  datastore.RawTransactionStore
		pkStore   datastore.PublicKeyStore
		networkID *big.Int
	}
	tests := []struct {
		name string
		args args
		want *Transaction
	}{
		{
			"success",
			args{
				datastoretest.NewMockTransactionStore(nil),
				datastoretest.NewMockRawTransactionStore(nil),
				datastoretest.NewMockPublicKeyStore(nil),
				big.NewInt(1),
			},
			&Transaction{
				datastoretest.NewMockTransactionStore(nil),
				datastoretest.NewMockRawTransactionStore(nil),
				datastoretest.NewMockPublicKeyStore(nil),
				big.NewInt(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionProcessor(tt.args.store, tt.args.rawStore, tt.args.pkStore, tt.args.networkID); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewTransactionProcessor() = %v, want %v", got, tt.want)
			}
		})
	}
}
