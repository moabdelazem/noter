package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// TxFn is a function that executes within a transaction
type TxFn func(pgx.Tx) error

// WithTransaction executes the given function within a transaction
func (db *DB) WithTransaction(ctx context.Context, fn TxFn) error {
	// Begin transaction
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Execute the function
	if err := fn(tx); err != nil {
		// If the function returns an error, rollback the transaction
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// If the function succeeds, commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
