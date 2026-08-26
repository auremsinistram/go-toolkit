package transaction

import (
	"context"

	"github.com/auremsinistram/go-errors"
	"github.com/jackc/pgx/v5"
)

func HandleTx(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return errors.Wrap(e, "transaction - HandleTx - #1")
		}
	} else {
		if e := tx.Commit(ctx); e != nil {
			return errors.Wrap(e, "transaction - HandleTx - #2")
		}
	}

	return err
}
