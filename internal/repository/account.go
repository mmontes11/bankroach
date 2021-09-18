package repository

import (
	"context"
	"database/sql"

	"github.com/mmontes11/bankroach/internal/model"
)

type accountRepo struct{}

type TransferParams struct {
	Source      model.AccountID
	Destination model.AccountID
	Amount      int64
}

type AccountRepo interface {
	Insert(ctx context.Context, tx *sql.Tx, balance int64) (model.AccountID, error)
	List(ctx context.Context, tx *sql.Tx) ([]model.Account, error)
	FindWithMinBalance(ctx context.Context, tx *sql.Tx, minBalance int64) (model.AccountID, error)
	Transfer(ctx context.Context, tx *sql.Tx, params TransferParams) error
}

func (r *accountRepo) Insert(
	ctx context.Context,
	tx *sql.Tx,
	balance int64,
) (model.AccountID, error) {
	sql := `INSERT INTO accounts(balance) VALUES($1) RETURNING id`
	row := tx.QueryRowContext(ctx, sql, balance)

	var id model.AccountID
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *accountRepo) List(ctx context.Context, tx *sql.Tx) ([]model.Account, error) {
	sql := `SELECT id, balance FROM accounts`

	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account
	for rows.Next() {
		var id model.AccountID
		var balance int64

		if err := rows.Scan(&id, &balance); err != nil {
			return nil, err
		}

		accounts = append(accounts, model.Account{
			ID:      id,
			Balance: balance,
		})
	}

	return accounts, nil
}

func (r *accountRepo) FindWithMinBalance(
	ctx context.Context,
	tx *sql.Tx,
	minBalance int64,
) (model.AccountID, error) {
	query := `SELECT id FROM accounts WHERE balance >= $1`
	row := tx.QueryRowContext(ctx, query, minBalance)

	var account model.AccountID
	if err := row.Scan(&account); err != nil {
		if err == sql.ErrNoRows {
			return 0, model.ErrNoAccountsWithBalance
		}
		return 0, err
	}

	return account, nil
}

func (r *accountRepo) Transfer(
	ctx context.Context,
	tx *sql.Tx,
	params TransferParams,
) error {
	sql := `UPDATE accounts SET balance = balance - $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, sql, params.Amount, params.Source)
	if err != nil {
		return err
	}

	sql = `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, sql, params.Amount, params.Destination)
	if err != nil {
		return err
	}

	return nil
}

func NewAccountRepo() AccountRepo {
	return &accountRepo{}
}
