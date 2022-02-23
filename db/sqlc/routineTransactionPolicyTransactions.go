package db

import "context"

type RoutineTransactionPolicyResult struct {
	RoutineTransactionPolicy RoutineTransactionPolicy `json:"routine_transaction_policy"`
}

func (store *SQLStore) AddRoutineTransactionPolicy(ctx context.Context, arg CreateRoutineTransactionPolicyParams) (RoutineTransactionPolicyResult, error) {
	var result RoutineTransactionPolicyResult
	err := store.WithTransaction(ctx, store.db, func(q *Queries) (txErr error) {

		result.RoutineTransactionPolicy, txErr = q.CreateRoutineTransactionPolicy(ctx, arg)
		if txErr != nil {
			return txErr
		}
		return nil
	})
	return result, err
}
