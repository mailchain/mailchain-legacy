package settings

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values/valuestest"
	"github.com/mailchain/mailchain/internal/protocols"
)

func Test_senders(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		s values.Store
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			"success",
			args{
				func() values.Store {
					m := valuestest.NewMockStore(mockCtrl)
					// m.EXPECT().IsSet("sentstore.kind")
					return m
				}(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := senders(tt.args.s)
			if (got == nil) != tt.wantNil {
				t.Errorf("senders() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestSenders_Produce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		clients map[string]SenderClient
	}
	type args struct {
		client string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			fields{
				map[string]SenderClient{
					"client": func() SenderClient {
						s := valuestest.NewMockStore(mockCtrl)
						s.EXPECT().IsSet("senders.ethereum-relay.base-url").Return(false)
						return relaySender(s, protocols.Ethereum)
					}(),
				},
			},
			args{
				"client",
			},
			false,
			false,
		},
		{
			"err-no-client",
			fields{
				map[string]SenderClient{
					"client": func() SenderClient {
						s := valuestest.NewMockStore(mockCtrl)
						return relaySender(s, protocols.Ethereum)
					}(),
				},
			},
			args{
				"invalid-client",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Senders{
				clients: tt.fields.clients,
			}
			got, err := s.Produce(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("Senders.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("Senders.Produce() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}
		})
	}
}
