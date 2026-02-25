package storage

import (
	"context"
	"errors"
	"fmt"
	"simple-banking-system/module/account/model"

	"github.com/jackc/pgx/v5"
)

func (s *sqlStore) CreateAccount(ctx context.Context, data *model.Account) (*model.Account, error) {
	if data == nil {
		return nil, model.ErrInvalidRequest
	}

	query := `
		INSERT INTO accounts (owner_name, balance)
		VALUES ($1, $2)
		RETURNING id, owner_name, balance, created_at
	`

	var account model.Account
	err := s.db.QueryRow(ctx, query, data.OwnerName, data.Balance).Scan(
		&account.ID,
		&account.OwnerName,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: create account: %v", model.ErrDBOperation, err)
	}

	return &account, nil
}

func (s *sqlStore) GetAccount(ctx context.Context, data *model.Account) (*model.Account, error) {
	query := `
		SELECT id, owner_name, balance, created_at
		FROM accounts
		WHERE id = $1
		LIMIT 1
	`
	var account model.Account
	err := s.db.QueryRow(ctx, query, data.OwnerName, data.Balance).Scan(
		&account.ID,
		&account.OwnerName,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrAccountNotFound
		}

		return nil, fmt.Errorf("%w: get account: %v", model.ErrDBOperation, err)
	}

	return &account, nil
}

func (s *sqlStore) ListAccounts(ctx context.Context, limit, offset int32) ([]model.Account, error) {
	query := `
		SELECT id, owner_name, balance, created_at
		FROM accounts
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: list accounts: %v", model.ErrDBOperation, err)
	}
	defer rows.Close()

	accounts := make([]model.Account, 0)
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(
			&acc.ID,
			&acc.OwnerName,
			&acc.Balance,
			&acc.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: scan account row: %v", model.ErrDBOperation, err)
		}
		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: rows error: %v", model.ErrDBOperation, err)
	}

	return accounts, nil
}
