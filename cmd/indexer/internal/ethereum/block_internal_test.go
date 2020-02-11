package ethereum

import (
	"testing"

	"github.com/mailchain/mailchain/cmd/indexer/internal/actions"
	"github.com/mailchain/mailchain/cmd/indexer/internal/actions/actionstest"
	"github.com/stretchr/testify/assert"
)

func TestNewBlockProcessor(t *testing.T) {
	type args struct {
		tx actions.Transaction
	}
	tests := []struct {
		name string
		args args
		want *Block
	}{
		{
			"success",
			args{
				actionstest.NewMockTransaction(nil),
			},
			&Block{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlockProcessor(tt.args.tx); !assert.IsType(t, tt.want, got) {
				t.Errorf("NewBlockProcessor() = %v, want %v", got, tt.want)
			}
		})
	}
}
