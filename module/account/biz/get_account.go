package biz

import (
	"context"
	"fmt"
	"simple-banking-system/module/account/model"
)

type GetAccountStore interface {
	GetAccount(ctx context.Context, id int64) (*model.Account, error)
}

type GetAccountBiz struct {
	store GetAccountStore
}

func NewGetAccountBiz(store GetAccountStore) *GetAccountBiz {
	return &GetAccountBiz{store: store}
}

func (biz *GetAccountBiz) GetAccount(ctx context.Context, id int64) (*model.Account, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: account id must be > 0", model.ErrInvalidRequest)
	}

	account, err := biz.store.GetAccount(ctx, id)

	if err != nil {
		return nil, err
	}

	return account, nil
}
