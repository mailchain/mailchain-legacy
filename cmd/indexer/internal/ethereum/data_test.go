package ethereum_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func getTx(t *testing.T, hash string) *types.Transaction {
	transactionJSON, err := ioutil.ReadFile(fmt.Sprintf("./testdata/transactions/%s.json", hash))
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil
	}

	tx := &types.Transaction{}
	if err := tx.UnmarshalJSON(transactionJSON); err != nil {
		assert.FailNow(t, err.Error())
		return nil
	}
	return tx
}
