package postgres

// import (
// 	"context"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgconn"
// )

// type DbPool interface {
// 	Close()
// 	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
// 	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
// 	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
// 	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
// }
