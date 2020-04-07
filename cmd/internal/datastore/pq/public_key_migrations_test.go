package pq

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type anyMatcher struct{}

func (a anyMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}
func TestMigratePublicKey(t *testing.T) {
	type args struct {
		db *sql.DB
		up bool
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"up",
			args{
				func() *sql.DB {
					db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(anyMatcher{}))
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}
					m.ExpectExec("*").WillReturnResult(sqlmock.NewResult(1, 1))
					m.ExpectQuery("*").WillReturnRows(
						sqlmock.NewRows([]string{"id", "applied_at"}).
							AddRow("1581972558643-create-public-key-table", time.Now()))

					return db
				}(),
				true,
			},
			0,
			false,
		},
		{
			"down",
			args{
				func() *sql.DB {
					db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(anyMatcher{}))
					if err != nil {
						t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
					}
					m.ExpectExec("*").WillReturnResult(sqlmock.NewResult(1, 1))
					m.ExpectQuery("*").WillReturnRows(
						sqlmock.NewRows([]string{"id", "applied_at"}).
							AddRow("1581972558643-create-public-key-table", time.Now()))

					return db
				}(),
				true,
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MigratePublicKey(tt.args.db, tt.args.up)
			if (err != nil) != tt.wantErr {
				t.Errorf("MigratePublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MigratePublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
