package settings

import (
	"testing"
)

func TestSentStoreMailchain_Produce(t *testing.T) {
	tests := []struct {
		name    string
		m       SentStoreMailchain
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			SentStoreMailchain{},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := SentStoreMailchain{}
			got, err := m.Produce()
			if (err != nil) != tt.wantErr {
				t.Errorf("SentStoreMailchain.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("SentStoreMailchain.Produce() nil = %v, wantNil %v", err, tt.wantNil)
				return
			}
		})
	}
}
