package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"lingua-evo/storage"
	"testing"
)

func TestDatabase_AddWord(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		w   *storage.Word
	}
	pool, _ := pgxpool.Connect(context.Background(), "postgres://postgres:5623@localhost:5432/postgres")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{name: "first_word", fields: fields{
			db: pool,
		}, args: args{ctx: context.Background(), w: &storage.Word{
			UserName:  "mak",
			Value:     "first",
			Language:  storage.Language{Origin: "Eng", Translate: "Rus"},
			Translate: []string{"первый"},
			Example:   []storage.Example{{Value: "1", Translate: "1"}},
		}}, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				db: tt.fields.db,
			}
			if err := d.AddWord(tt.args.ctx, tt.args.w); err != tt.wantErr {
				t.Errorf("AddWord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
