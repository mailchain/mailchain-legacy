package substrate_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

func getBlock(t *testing.T, hash string) *types.Block {
	blockJSON, err := ioutil.ReadFile(fmt.Sprintf("./testdata/blocks/%s.json", hash))
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil
	}

	type jsonError struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
	// A value of this type can a JSON-RPC request, notification, successful response or
	// error response. Which one it is depends on the fields.
	type jsonrpcMessage struct {
		Version string          `json:"jsonrpc,omitempty"`
		ID      json.RawMessage `json:"id,omitempty"`
		Method  string          `json:"method,omitempty"`
		Params  json.RawMessage `json:"params,omitempty"`
		Error   *jsonError      `json:"error,omitempty"`
		Result  json.RawMessage `json:"result,omitempty"`
	}

	msg := &jsonrpcMessage{}
	if err := json.Unmarshal(blockJSON, &msg); err != nil {
		assert.FailNow(t, err.Error())
		return nil
	}

	var result types.SignedBlock
	if err := json.Unmarshal(msg.Result, &result); err != nil {
		assert.FailNow(t, err.Error())
		return nil
	}

	return &result.Block
}
