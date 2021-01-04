package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mailchain/mailchain/cmd/internal/datastore"
	"github.com/mailchain/mailchain/cmd/internal/http/params"
	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/pkg/errors"
)

// HandleToRequest accepts all relay requests and routes then to the new URL as required.
func HandleToRequest(ts datastore.TransactionStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req, err := parseGetEnvelopesRequest(r)
		if err != nil {
			errs.JSONWriter(w, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		txs, err := ts.GetTransactionsTo(ctx, req.Protocol, req.Network, req.addressBytes)
		if err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		envelopes := []Envelope{}
		for x := range txs {
			envelopes = append(envelopes,
				Envelope{
					From:        encoding.EncodeHexZeroX(txs[x].From),
					To:          encoding.EncodeHexZeroX(txs[x].To),
					BlockHash:   encoding.EncodeHexZeroX(txs[x].BlockHash),
					BlockNumber: fmt.Sprint(txs[x].BlockNumber),
					Hash:        encoding.EncodeHexZeroX(txs[x].Hash),
					Data:        encoding.EncodeHexZeroX(txs[x].Data),

					Value:    txs[x].Value.String(),
					GasPrice: txs[x].GasPrice.String(),
					GasUsed:  txs[x].GasUsed.String(),
				},
			)
		}

		if err := json.NewEncoder(w).Encode(getEnvelopesResponse{Envelopes: envelopes}); err != nil {
			errs.JSONWriter(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

// GetEnvelopesRequest get mailchain envelopes
// swagger:parameters GetEnvelopes
type GetEnvelopesRequest struct {
	// Address to use when looking for envelopes.
	//
	// in: query
	// required: true
	// example: 0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae
	// pattern: 0x[a-fA-F0-9]{40}
	Address string `json:"address"`

	// Network to use when looking for envelopes.
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: query
	// required: true
	// example: goerli
	Network string `json:"network"`

	// Protocol to use when looking for envelopes.
	//
	// enum: ethereum
	// in: query
	// required: true
	// example: ethereum
	Protocol string `json:"protocol"`

	addressBytes []byte
}

func parseGetEnvelopesRequest(r *http.Request) (*GetEnvelopesRequest, error) {
	protocol, err := params.QueryRequireProtocol(r)
	if err != nil {
		return nil, err
	}

	network, err := params.QueryRequireNetwork(r)
	if err != nil {
		return nil, err
	}

	addr, err := params.QueryRequireAddress(r)
	if err != nil {
		return nil, err
	}

	addressBytes, err := address.DecodeByProtocol(addr, protocol)
	if err != nil {
		return nil, err
	}

	return &GetEnvelopesRequest{
		Address:      addr,
		addressBytes: addressBytes,
		Network:      network,
		Protocol:     protocol,
	}, nil
}

// getEnvelopesResponse Holds the response messages
//
// swagger:response GetEnvelopesResponse
type getEnvelopesResponse struct {
	// in: body
	Envelopes []Envelope `json:"envelopes"`
}

type Envelope struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Data        string `json:"data"`
	BlockHash   string `json:"block-hash"`
	BlockNumber string `json:"block-number"`
	Hash        string `json:"hash"`

	Value    string `json:"value"`
	GasUsed  string `json:"gas-used"`
	GasPrice string `json:"gas-price"`
}
