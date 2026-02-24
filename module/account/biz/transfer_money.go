package biz

import (
	"context"
	"fmt"

	"simple-banking-system/module/account/model"
)

type TransferMoneyResult struct {
	Transfer    *model.Transfer `json:"transfer"`
	FromAccount *model.Account  `json:"from_account"`
	ToAccount   *model.Account  `json:"to_account"`
	FromEntry   *model.Entry    `json:"from_entry"`
	ToEntry     *model.Entry    `json:"to_entry"`
}

type TransferMoneyStore interface {
	TransferMoney(
		ctx context.Context,
		fromAccountID int64,
		toAccountID int64,
		amount int64,
	) (*TransferMoneyResult, error)
}

type transferMoneyBiz struct {
	store TransferMoneyStore
}

func NewTransferMoneyBiz(store TransferMoneyStore) *transferMoneyBiz {
	return &transferMoneyBiz{store: store}
}

func (biz *transferMoneyBiz) TransferMoney(
	ctx context.Context,
	req *model.CreateTransferRequest,
) (*TransferMoneyResult, error) {
	if req == nil {
		return nil, model.ErrInvalidRequest
	}

	if req.FromAccountID <= 0 {
		return nil, fmt.Errorf("%w: from_account_id must be > 0", model.ErrInvalidRequest)
	}

	if req.ToAccountID <= 0 {
		return nil, fmt.Errorf("%w: to_account_id must be > 0", model.ErrInvalidRequest)
	}

	if req.FromAccountID == req.ToAccountID {
		return nil, model.ErrSameAccountTransfer
	}

	if req.Amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be > 0", model.ErrInvalidRequest)
	}

	result, err := biz.store.TransferMoney(
		ctx,
		req.FromAccountID,
		req.ToAccountID,
		req.Amount,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
