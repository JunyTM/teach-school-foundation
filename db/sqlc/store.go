package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provides all func to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferResult, error)
}

// store provides all func to execute db queries and transactions
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

// execTx executes a function within a database transaction
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil) // Parameter 2 to set Options allowing for Isolation Level, ReadOnly, etc.
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx_error: %v - tx.Rollback() error: %v", err, rbErr)
		} // Rollback if error
		return err
	}

	return tx.Commit() // Commit if no error
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// var txKey = struct{}{}

// Transfer preforms a mone transfer from one account to the other
func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferResult, error) {
	var result TransferResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		// txName := ctx.Value(txKey)

		// fmt.Println(txName, "crate transfer")
		result.Transfer, err = q.CreateTranfer(ctx, CreateTranferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
			Currency:      "USD",
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "created entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{ // Create Entry for the From Account
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "created entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{ // Create Entry for the To Account
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, s, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, s, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	s *SQLStore, fromAccountID int64,
	fromAmount int64,
	toAccountID int64,
	toAmount int64,
) (fromAccount Account, toAccount Account, err error) {
	fromAccount, err = s.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountID,
		Amount: fromAmount,
	})
	if err != nil {
		return
	}

	toAccount, err = s.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountID,
		Amount: toAmount,
	})
	return
}
