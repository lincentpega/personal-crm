package txcontext

import (
	"context"
	"database/sql"
)

type ctxKey int

const txKey ctxKey = 0

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}
