package relay

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/internal/mailbox/signer"
	"github.com/mailchain/mailchain/internal/testutil"
	"github.com/mailchain/mailchain/sender"
	"github.com/mailchain/mailchain/sender/sendertest"
	"github.com/stretchr/testify/assert"
)

func Test_createAddress(t *testing.T) {
	type args struct {
		baseURL string
		chain   string
		network string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"success-trailing",
			args{
				"https://relay.mailchain.xyz/",
				"ethereum",
				"mainnet",
			},
			"https://relay.mailchain.xyz/json-rpc/ethereum/mainnet",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createAddress(tt.args.baseURL, tt.args.chain, tt.args.network); got != tt.want {
				t.Errorf("createAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(nil)
	type args struct {
		baseURL string
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantNumSenders int
	}{
		{
			"success",
			args{
				server.URL,
			},
			false,
			5,
		},
		{
			"err-client",
			args{
				"ttps://bad#host",
			},
			true,
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !assert.Len(got.senders, tt.wantNumSenders) {
					t.Errorf("NewClient().senders = %v, want %v", got, tt.wantNumSenders)
				}
			}

		})
	}
}

func TestClient_Send(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		senders map[string]sender.Message
	}
	type args struct {
		ctx     context.Context
		network string
		to      []byte
		from    []byte
		data    []byte
		signer  signer.Signer
		opts    sender.MessageOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"success",
			fields{
				map[string]sender.Message{
					"mainnet": func() sender.Message {
						m := sendertest.NewMockMessage(mockCtrl)
						m.EXPECT().Send(
							context.Background(),
							"mainnet",
							testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
							testutil.MustHexDecodeString("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
							[]byte("transactionDataValue"),
							nil,
							nil,
						).Return(nil)
						return m
					}(),
				},
			},
			args{
				context.Background(),
				"mainnet",
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				testutil.MustHexDecodeString("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
				[]byte("transactionDataValue"),
				nil,
				nil,
			},
			false,
		},
		{
			"success",
			fields{
				map[string]sender.Message{},
			},
			args{
				context.Background(),
				"mainnet",
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				testutil.MustHexDecodeString("4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2"),
				[]byte("transactionDataValue"),
				nil,
				nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				senders: tt.fields.senders,
			}
			if err := c.Send(tt.args.ctx, tt.args.network, tt.args.to, tt.args.from, tt.args.data, tt.args.signer, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Client.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
