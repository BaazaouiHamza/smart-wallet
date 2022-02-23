package db

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/multierr"
)

type Store interface {
	Querier
	AddRoutineTransactionPolicy(ctx context.Context, arg CreateRoutineTransactionPolicyParams) (RoutineTransactionPolicyResult, error)
}

//SQLStoreprovides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) WithTransaction(
	ctx context.Context, db *sql.DB, transaction func(qa *Queries) (txErr error),
) (err error) {
	var tx *sql.Tx

	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return fmt.Errorf("could not start transcation: %w", err)
	}

	defer func() {
		v := recover()

		switch {
		case v != nil:
			_ = tx.Rollback()

			panic(v)
		case err != nil:
			if rbErr := tx.Rollback(); rbErr != nil {
				err = multierr.Combine(
					err,
					fmt.Errorf("could not rollback transaction: %w", rbErr),
				)
			}
		default:
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("could not commit transaction: %w", err)
			}
		}
	}()

	err = transaction(&Queries{tx})

	return
}
