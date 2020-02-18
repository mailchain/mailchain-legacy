package pq

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	type args struct {
		user         string
		password     string
		databaseName string
		host         string
		sslmode      string
		port         int
	}
	tests := []struct {
		name    string
		args    args
		want    *sqlx.DB
		wantErr bool
	}{
		{
			"err",
			args{
				"user",
				"password",
				"database",
				"localhost",
				"disable",
				99999999,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnection(tt.args.user, tt.args.password, tt.args.databaseName, tt.args.host, tt.args.sslmode, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("NewConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
