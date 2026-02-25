package storage

import (
	"context"
	"errors"
	"fmt"
	"simple-banking-system/module/account/biz"
	"simple-banking-system/module/account/model"

	"github.com/jackc/pgx/v5"
)

func (s *sqlStore) GetAccountForUpdateTx(ctx context.Context, tx pgx.Tx, id int64) (*model.Account, error) {
	query := `
		SELECT id, owner_name, balance, created_at
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`
	var acc model.Account
	err := tx.QueryRow(ctx, query, id).Scan(
		&acc.ID,
		&acc.OwnerName,
		&acc.Balance,
		&acc.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrAccountNotFound
		}

		return nil, fmt.Errorf("%w: lock account: %v", model.ErrDBOperation, err)
	}

	return &acc, nil
}

func (s *sqlStore) CreateTransferTx(ctx context.Context, tx pgx.Tx, fromAccountID int64, toAccountID int64, amount int64) (*model.Transfer, error) {
	query := `
		INSERT INTO transfers (from_account_id, to_account_id, amount)
		VALUES ($1, $2, $3)
		RETURNING id, from_account_id, to_account_id, amount, created_at
	`
	var t model.Transfer
	err := tx.QueryRow(ctx, query, fromAccountID, toAccountID, amount).Scan(
		&t.ID,
		&t.FromAccountID,
		&t.ToAccountID,
		&t.Amount,
		&t.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: create transfer: %v", model.ErrDBOperation, err)
	}

	return &t, nil
}

func (s *sqlStore) CreateEntryTx(ctx context.Context, tx pgx.Tx, accountID int64, transferID int64, amount int64) (*model.Entry, error) {
	query := `
		INSERT INTO entries (from_account_id, to_account_id, amount)
		VALUES ($1, $2, $3)
		RETURNING id, from_account_id, to_account_id, amount, created_at
	`
	var e model.Entry
	err := tx.QueryRow(ctx, query, accountID, transferID, amount).Scan(
		&e.ID,
		&e.AccountID,
		&e.TransferID,
		&e.Amount,
		&e.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: create entry: %v", model.ErrDBOperation, err)
	}

	return &e, nil
}

func (s *sqlStore) addAccountBalanceTx(ctx context.Context, tx pgx.Tx, accountID int64, delta int64) (*model.Account, error) {
	query := `
		UPDATE accounts
		SET balance = balance + $2
		WHERE id = $1
		RETURNING id, owner_name, balance, created_at
	`
	var acc model.Account
	err := tx.QueryRow(ctx, query, accountID, delta).Scan(
		&acc.ID,
		&acc.OwnerName,
		&acc.Balance,
		&acc.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrAccountNotFound
		}

		return nil, fmt.Errorf("%w: update account balance: %v", model.ErrDBOperation, err)
	}

	return &acc, nil
}

func (s *sqlStore) TransferMoney(ctx context.Context, fromAccountID int64, toAccountID int64, amount int64) (*biz.TransferMoneyResult, error) {
	var result biz.TransferMoneyResult
	err := s.execTx(ctx, func(tx pgx.Tx) error {
		firstID, secondID := fromAccountID, toAccountID

		if firstID > secondID {
			firstID, secondID = secondID, firstID
		}

		acc1, err := s.GetAccountForUpdateTx(ctx, tx, firstID)

		if err != nil {
			return err
		}

		acc2, err := s.GetAccountForUpdateTx(ctx, tx, secondID)

		if err != nil {
			return err
		}

		var fromAcc, toAcc *model.Account

		if acc1.ID == fromAccountID {
			fromAcc = acc1
			toAcc = acc2
		} else {
			fromAcc = acc2
			toAcc = acc1
		}

		if fromAcc.Balance < amount {
			return model.ErrInsufficientFunds
		}

		transfer, err := s.CreateTransferTx(ctx, tx, fromAccountID, toAccountID, amount)

		if err != nil {
			return err
		}

		fromEntry, err := s.CreateEntryTx(ctx, tx, fromAccountID, transfer.ID, -amount)

		if err != nil {
			return err
		}

		toEntry, err := s.CreateEntryTx(ctx, tx, toAccountID, transfer.ID, amount)

		if err != nil {
			return err
		}

		var updateFirst, updateSecond *model.Account

		if firstID == fromAccountID {
			updateFirst, err = s.addAccountBalanceTx(ctx, tx, fromAccountID, -amount)

			if err != nil {
				return err
			}

			updateSecond, err = s.addAccountBalanceTx(ctx, tx, toAccountID, amount)

			if err != nil {
				return err
			}

			fromAcc = updateFirst
			toAcc = updateSecond
		} else {
			updateFirst, err = s.addAccountBalanceTx(ctx, tx, toAccountID, amount)

			if err != nil {
				return err
			}

			updateSecond, err = s.addAccountBalanceTx(ctx, tx, fromAccountID, -amount)

			if err != nil {
				return err
			}

			toAcc = updateFirst
			fromAcc = updateSecond
		}

		result = biz.TransferMoneyResult{
			Transfer:    transfer,
			FromAccount: fromAcc,
			ToAccount:   toAcc,
			FromEntry:   fromEntry,
			ToEntry:     toEntry,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}
