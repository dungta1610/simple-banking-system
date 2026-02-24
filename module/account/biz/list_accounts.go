package biz

import (
	"context"
	"fmt"

	"simple-banking-system/module/account/model"
)

type ListAccountsStore interface {
	ListAccounts(ctx context.Context, limit, offset int32) ([]model.Account, error)
}

type listAccountsBiz struct {
	store ListAccountsStore
}

func NewListAccountsBiz(store ListAccountsStore) *listAccountsBiz {
	return &listAccountsBiz{store: store}
}

func (biz *listAccountsBiz) ListAccounts(
	ctx context.Context,
	query *model.ListAccountsQuery,
) ([]model.Account, error) {
	if query == nil {
		return nil, model.ErrInvalidRequest
	}

	if err := query.Normalize(); err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInvalidRequest, err)
	}

	accounts, err := biz.store.ListAccounts(ctx, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
