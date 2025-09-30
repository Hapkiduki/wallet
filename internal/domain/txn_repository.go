package domain

import "context"

// TxnRepository defines the contract for transaction management.
type TxnRepository interface {
	// WithTransaction executes the given function within a database transaction.
	// If the function returns an error, the transaction is rolled back.
	// Otherwise, the transaction is committed.
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
