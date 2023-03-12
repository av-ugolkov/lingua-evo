package repository

import (
	"context"
	"database/sql"
	"testing"
)

func TestDatabase_AddWord(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		w   *Word
	}
	pool, _ := NewDB("postgres://postgres:5623@localhost:5432/postgres")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{name: "first_word", fields: fields{
			db: pool,
		}, args: args{ctx: context.Background(), w: &Word{
			Text:     "first",
			Language: "en_EN",
		}}, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				db: tt.fields.db,
			}
			if uuid, err := d.AddWord(tt.args.ctx, tt.args.w); err != tt.wantErr {
				t.Errorf("AddWord() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("AddWord() word uuid: %s", uuid)
			}
		})
	}
}
