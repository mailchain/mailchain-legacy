package etherscan

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFromResultHash(t *testing.T) {
	assert := assert.New(t)

	mockTxResult := []txResult{
		{From: "address1", Hash: "aaa111"},
		{From: "address2", Hash: "bbb222"},
		{From: "address3", Hash: "ccc333"},
	}

	type args struct {
		address string
		txList  *txList
	}

	testCases := []struct {
		name    string
		args    args
		want    common.Hash
		wantErr bool
	}{
		{
			"err-empty-transaction-list",
			args{
				"should-not-matter",
				&txList{
					Status:  "",
					Message: "",
					Result:  []txResult{},
				},
			},
			common.Hash{},
			true,
		},
		{
			"match-transaction-1",
			args{
				"address1",
				&txList{
					Status:  "",
					Message: "",
					Result:  mockTxResult,
				},
			},
			common.HexToHash("aaa111"),
			false,
		},
		{
			"match-transaction-3",
			args{
				"address3",
				&txList{
					Status:  "",
					Message: "",
					Result:  mockTxResult,
				},
			},
			common.HexToHash("ccc333"),
			false,
		},
		{
			"err-no-matching-transactions",
			args{
				"address11",
				&txList{
					Status:  "",
					Message: "",
					Result:  mockTxResult,
				},
			},
			common.Hash{},
			true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			hash, err := getFromResultHash(testCase.args.address, testCase.args.txList)

			if (err != nil) != testCase.wantErr {
				t.Errorf("getFromResultHash() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if !assert.Equal(testCase.want, hash) {
				t.Errorf("getFromResultHash() = %v, want %v", hash, testCase.want)
			}
		})
	}
}
