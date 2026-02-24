package biz

import (
	"context"
	"fmt"
	"simple-banking-system/module/account/model"
)

type CreateAccountStore interface {
	CreateAccount(ctx context.Context, data *model.Account) (*model.Account, error)
}

type CreateAccountBiz struct {
	store CreateAccountStore
}

func NewCreateAccountBiz(store CreateAccountStore) *CreateAccountBiz {
	return &CreateAccountBiz{store: store}
}

func (biz *CreateAccountBiz) CreateAccount(ctx context.Context, req *model.CreateAccountRequest) (*model.Account, error) {
	if req == nil {
		return nil, model.ErrInvalidRequest
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInvalidRequest, err)
	}

	account := &model.Account{
		OwnerName: req.OwnerName,
		Balance:   req.Balance,
	}

	created, err := biz.store.CreateAccount(ctx, account)

	if err != nil {
		return nil, err
	}

	return created, nil
}
